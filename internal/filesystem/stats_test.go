package filesystem_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/davidsbond/homelab/internal/filesystem"
)

func TestGetSummary(t *testing.T) {
	t.Parallel()

	summary, err := filesystem.GetSummary("/")
	assert.NoError(t, err)
	assert.NotZero(t, summary)
	assert.EqualValues(t, summary.Total, summary.Used+summary.Available)
	assert.EqualValues(t, 100, summary.PercentageUsed+summary.PercentageAvailable)
}
