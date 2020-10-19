package main

import (
	"context"
	"fmt"
	"os"

	"pkg.dsb.dev/app"
	"pkg.dsb.dev/flag"
	"pkg.dsb.dev/metrics"

	"github.com/davidsbond/homelab/internal/weather"
)

func main() {
	a := app.New(
		app.WithRunner(run),
		app.WithFlags(
			&flag.String{
				Name:        "api-key",
				Usage:       "API key for authentication with weatherapi.com",
				EnvVar:      "API_KEY",
				Required:    true,
				Destination: &apiKey,
			},
			&flag.String{
				Name:        "location",
				Usage:       "The location to query weather information for",
				EnvVar:      "LOCATION",
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
	apiKey   string
	location string
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

	return metrics.Push()
}
