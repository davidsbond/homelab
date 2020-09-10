package pihole_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/davidsbond/homelab/internal/pihole"
)

func TestPiHole_Summary(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	ph, err := pihole.New("http://127.0.0.1:80")
	assert.NoError(t, err)

	summary, err := ph.Summary(ctx)
	assert.NoError(t, err)
	assert.NotZero(t, summary)
}
