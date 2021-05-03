package alertmanager_test

import (
	"context"

	"github.com/davidsbond/homelab/internal/alertmanager"
)

type (
	MockDispatcher struct {
		webhook alertmanager.Webhook
	}
)

func (m *MockDispatcher) Dispatch(ctx context.Context, webhook alertmanager.Webhook) error {
	m.webhook = webhook
	return nil
}
