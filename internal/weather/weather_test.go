package weather_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/davidsbond/homelab/internal/weather"
)

func TestClient_GetWeather(t *testing.T) {
	t.Parallel()

	w, err := weather.New("")

	assert.NoError(t, err)
	assert.NotNil(t, w)

	ctx := context.Background()
	result, err := w.GetWeather(ctx, "London")

	assert.NoError(t, err)
	assert.NotNil(t, result)
}
