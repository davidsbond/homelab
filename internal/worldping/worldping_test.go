package worldping_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/davidsbond/homelab/internal/worldping"
)

func TestRun(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	stats, err := worldping.Run(ctx, false)

	assert.NoError(t, err)
	assert.NotNil(t, stats)
}
