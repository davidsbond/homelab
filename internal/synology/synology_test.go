package synology_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/davidsbond/homelab/internal/synology"
)

func TestClient_SystemInfo(t *testing.T) {
	ctx := context.Background()

	cl, err := synology.New("", "", "")
	assert.NoError(t, err)
	assert.NotNil(t, cl)

	info, err := cl.SystemInfo(ctx)
	assert.NotNil(t, info)
	assert.NoError(t, err)
}
