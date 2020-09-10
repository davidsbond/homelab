package speedtest_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/davidsbond/homelab/internal/speedtest"
)

func TestTester_Test(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	result, err := speedtest.New().Test(ctx)
	assert.NoError(t, err)
	assert.NotZero(t, result.Latency)
}
