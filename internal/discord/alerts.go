// Package discord provide types for interacting with discord webhooks.
package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/davidsbond/homelab/internal/alertmanager"
)

type (
	// The AlertDispatcher type is used to convert alert manager alerts into discord webhooks.
	AlertDispatcher struct {
		webhookURL string
		client     *http.Client
	}

	// The Webhook type describes the payload of a discord webhook.
	Webhook struct {
		Content string  `json:"content"`
		Embeds  []Embed `json:"embeds"`
	}

	// The Embed type represents embeddable data in a discord message.
	Embed struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Colour      int       `json:"color"`
		Fields      []Field   `json:"fields"`
		URL         string    `json:"url"`
		Timestamp   time.Time `json:"timestamp"`
	}

	// The Field type represents a single field in a discord message.
	Field struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	// The APIError type represents an error returned from the discord API.
	APIError struct {
		Message string `json:"message"`
	}
)

// Constants for colours.
const (
	ColourRed   = 0x992D22
	ColourGreen = 0x2ECC71
	ColourGrey  = 0x95A5A6
)

// NewAlertDispatcher returns a new AlertDispatcher that will convert alert manager webhooks into discord webhooks
// and send them to the given URL.
func NewAlertDispatcher(webhookURL string) *AlertDispatcher {
	return &AlertDispatcher{
		webhookURL: webhookURL,
		client: &http.Client{
			Timeout: time.Minute,
		},
	}
}

// Dispatch an alert manager webhook to a discord webhook.
func (a *AlertDispatcher) Dispatch(ctx context.Context, webhook alertmanager.Webhook) error {
	for _, alert := range webhook.Alerts {
		var msg Webhook

		embed := Embed{
			Title:       webhook.CommonLabels["alertname"],
			Description: webhook.CommonAnnotations["summary"],
			Colour:      ColourGrey,
			URL:         alert.GeneratorURL,
			Timestamp:   alert.StartsAt,
		}

		switch alert.Status {
		case alertmanager.StatusFiring:
			embed.Colour = ColourRed
		case alertmanager.StatusResolved:
			embed.Colour = ColourGreen
		}

		if embed.Description != "" {
			msg.Content = embed.Description
		}

		for k, v := range alert.Labels {
			embed.Fields = append(embed.Fields, Field{
				Name:  k,
				Value: v,
			})
		}

		for k, v := range alert.Annotations {
			embed.Fields = append(embed.Fields, Field{
				Name:  k,
				Value: v,
			})
		}

		msg.Embeds = append(msg.Embeds, embed)

		body, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.webhookURL, bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := a.client.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode < http.StatusIMUsed {
			continue
		}

		var e APIError
		if err = json.NewDecoder(resp.Body).Decode(&e); err != nil {
			return err
		}
		return e
	}

	return nil
}

func (err APIError) Error() string {
	return err.Message
}
