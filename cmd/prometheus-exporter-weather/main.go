package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/davidsbond/homelab/internal/weather"

	"github.com/urfave/cli/v2"
	"pkg.dsb.dev/app"
	"pkg.dsb.dev/cron"
	"pkg.dsb.dev/metrics"
)

const (
	defaultFrequency = time.Hour / 3
)

func main() {
	a := app.New(
		app.WithRunner(run),
		app.WithFlags(
			&cli.DurationFlag{
				Name:        "frequency",
				Usage:       "How often to query weather data",
				EnvVars:     []string{"FREQUENCY"},
				Value:       defaultFrequency,
				Destination: &frequency,
			},
			&cli.StringFlag{
				Name:        "api-key",
				Usage:       "API key for authentication with weatherapi.com",
				EnvVars:     []string{"API_KEY"},
				Required:    true,
				Destination: &apiKey,
			},
			&cli.StringFlag{
				Name:        "location",
				Usage:       "The location to query weather information for",
				EnvVars:     []string{"LOCATION"},
				Required:    true,
				Destination: &location,
			},
		),
	)

	if err := a.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

var (
	frequency time.Duration
	apiKey    string
	location  string
)

func run(ctx context.Context) error {
	metrics.Register(
		temperature, windSpeed, windDirection, uvIndex,
		gustSpeed, pressure, precipitation, humidity,
	)

	client, err := weather.New(apiKey)
	if err != nil {
		return err
	}

	return cron.Every(ctx, frequency, func(ctx context.Context) error {
		results, err := client.GetWeather(ctx, location)
		if err != nil {
			return err
		}

		temperature.WithLabelValues(location).Set(results.TempC)
		windSpeed.WithLabelValues(location).Set(results.WindMph)
		windDirection.WithLabelValues(location).Set(results.WindDegree)
		uvIndex.WithLabelValues(location).Set(results.Uv)
		gustSpeed.WithLabelValues(location).Set(results.GustMph)
		pressure.WithLabelValues(location).Set(results.PressureMb)
		precipitation.WithLabelValues(location).Set(results.PrecipMm)
		humidity.WithLabelValues(location).Set(results.Humidity)
		return nil
	})
}
