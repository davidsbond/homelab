package coronavirus_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/davidsbond/homelab/internal/coronavirus"
)

func TestClient_GetNewCaseCount(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cl, err := coronavirus.New()
	assert.NoError(t, err)
	assert.NotNil(t, cl)

	count, err := cl.GetSummary(ctx, time.Now())
	assert.NoError(t, err)
	assert.NotZero(t, count)
}
