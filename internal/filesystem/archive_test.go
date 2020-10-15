package filesystem_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/davidsbond/homelab/internal/filesystem"
)

func TestArchive(t *testing.T) {
	ctx := context.Background()

	output, err := os.Create("output.tar.gz")
	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, output.Close())
		assert.NoError(t, os.Remove("output.tar.gz"))
	})

	assert.NoError(t, filesystem.Archive(ctx, "test", output))
}
