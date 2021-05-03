// Package alertmanager provides types for working with alert manager webhooks.
package alertmanager

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"pkg.dsb.dev/transport"
)

type (
	// The Webhook type represents the structure of an alert manager webhook.
	Webhook struct {
		Alerts            []Alert           `json:"alerts"`
		CommonAnnotations map[string]string `json:"commonAnnotations"`
		CommonLabels      map[string]string `json:"commonLabels"`
		ExternalURL       string            `json:"externalURL"`
		GroupKey          string            `json:"groupKey"`
		GroupLabels       map[string]string `json:"groupLabels"`
		Receiver          string            `json:"receiver"`
		Status            string            `json:"status"`
	}

	// The Alert type represents a single alert within a webhook payload.
	Alert struct {
		Annotations  map[string]string `json:"annotations"`
		EndsAt       time.Time         `json:"endsAt"`
		GeneratorURL string            `json:"generatorURL"`
		Labels       map[string]string `json:"labels"`
		StartsAt     time.Time         `json:"startsAt"`
		Status       string            `json:"status"`
	}

	// The Dispatcher interface describes types that can handle alert manager webhooks and forward them somewhere.
	Dispatcher interface {
		Dispatch(ctx context.Context, webhook Webhook) error
	}

	// The WebhookHandler type is used to handle inbound alert manager webhooks via HTTP.
	WebhookHandler struct {
		transport.HTTP

		dispatcher Dispatcher
	}
)

// Constants for alert statuses.
const (
	StatusFiring   = "firing"
	StatusResolved = "resolved"
)

// NewWebhookHandler returns a new WebhookHandler instance that will forward webhooks to the given Dispatcher
// implementation.
func NewWebhookHandler(dispatcher Dispatcher) *WebhookHandler {
	return &WebhookHandler{
		dispatcher: dispatcher,
	}
}

// Register endpoints on the provided router.
func (wh *WebhookHandler) Register(r *mux.Router) {
	r.HandleFunc("/webhook", wh.Handle).Methods(http.MethodPost)
}

// Handle an inbound HTTP request.
func (wh *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var webhook Webhook
	if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
		wh.Error(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	if err := wh.dispatcher.Dispatch(r.Context(), webhook); err != nil {
		wh.Error(r.Context(), w, http.StatusInternalServerError, err.Error())
	}
}
