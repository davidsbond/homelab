package discord_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/davidsbond/homelab/internal/alertmanager"
	"github.com/davidsbond/homelab/internal/discord"
)

func TestAlertDispatcher_Dispatch(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name   string
		Input  alertmanager.Webhook
		Output discord.Webhook
	}{
		{
			Name: "It should convert an alertmanager webhook to a discord one",
			Input: alertmanager.Webhook{
				Alerts: []alertmanager.Alert{
					{
						Status: alertmanager.StatusFiring,
						Annotations: map[string]string{
							"description": "some description",
						},
						StartsAt:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						GeneratorURL: "https://example.com:9090",
					},
				},
				CommonAnnotations: map[string]string{
					"summary": "A man has fallen into the river in lego city",
				},
				CommonLabels: map[string]string{
					"alertname": "Alert!",
				},
				ExternalURL: "https://example.com:9093",
				GroupKey:    "{}:{alertname=\"Test\", job=\"prometheus24\"}",
				GroupLabels: map[string]string{
					"alertname": "Test",
					"job":       "prometheus24",
				},
				Receiver: "webhook",
				Status:   alertmanager.StatusFiring,
			},
			Output: discord.Webhook{
				Content: "A man has fallen into the river in lego city",
				Embeds: []discord.Embed{
					{
						Title:       "Alert!",
						Description: "A man has fallen into the river in lego city",
						Colour:      discord.ColourRed,
						URL:         "https://example.com:9090",
						Timestamp:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						Fields: []discord.Field{
							{
								Name:  "description",
								Value: "some description",
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go httpBin(ctx, t, tc.Output)

			dispatcher := discord.NewAlertDispatcher("http://localhost:8080")
			assert.NoError(t, dispatcher.Dispatch(ctx, tc.Input))
		})
	}
}

func httpBin(ctx context.Context, t *testing.T, expect discord.Webhook) {
	t.Helper()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var result discord.Webhook
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&result))
		assert.EqualValues(t, expect, result)
	})

	svr := &http.Server{Handler: handler, Addr: ":8080"}
	go func() {
		err := svr.ListenAndServe()
		switch {
		case errors.Is(err, http.ErrServerClosed):
			return
		default:
			assert.NoError(t, err)
		}
	}()

	<-ctx.Done()
	assert.NoError(t, svr.Shutdown(context.Background()))
}
