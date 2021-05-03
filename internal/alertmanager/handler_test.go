package alertmanager_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/davidsbond/homelab/internal/alertmanager"
)

func TestWebhookHandler_Handle(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name    string
		Webhook alertmanager.Webhook
	}{
		{
			Name: "It should forward a webhook to the dispatcher",
			Webhook: alertmanager.Webhook{
				Alerts: []alertmanager.Alert{
					{
						Status: alertmanager.StatusFiring,
						Labels: map[string]string{
							"alertname": "Test",
							"dc":        "eu-west-1",
							"instance":  "localhost:9090",
							"job":       "prometheus24",
						},
						Annotations: map[string]string{
							"description": "some description",
						},
						StartsAt:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						GeneratorURL: "https://example.com:9090/graph?g0.expr=go_memstats_alloc_bytes+%3E+0\u0026g0.tab=1",
					},
				},
				CommonAnnotations: map[string]string{
					"alertname": "Test",
					"dc":        "eu-west-1",
					"instance":  "localhost:9090",
					"job":       "prometheus24",
				},
				CommonLabels: map[string]string{
					"description": "some description",
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
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			dispatcher := &MockDispatcher{}
			handler := mux.NewRouter()
			alertmanager.NewWebhookHandler(dispatcher).Register(handler)

			body, err := json.Marshal(tc.Webhook)
			assert.NoError(t, err)

			r := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, r)
			assert.EqualValues(t, tc.Webhook, dispatcher.webhook)
		})
	}
}
