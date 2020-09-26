package homehub_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/davidsbond/homelab/internal/homehub"
)

func TestClient_Summary(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	cl, err := homehub.New("")
	assert.NoError(t, err)
	assert.NotNil(t, cl)

	summary, err := cl.Summary(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, summary)
}
