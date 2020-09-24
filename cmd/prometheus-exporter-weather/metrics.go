package main

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "weather"
	subsystem = "stats"
)

var (
	temperature = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "temperature",
		Help:      "Temperature in Celsius",
	}, []string{"location"})

	windSpeed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "wind_speed",
		Help:      "Wind speed in miles-per-hour",
	}, []string{"location"})

	windDirection = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "wind_direction",
		Help:      "Wind direction in degrees",
	}, []string{"location"})

	uvIndex = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "uv_index",
		Help:      "UV Index",
	}, []string{"location"})

	gustSpeed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "gust_speed",
		Help:      "Wind gust in miles per hour",
	}, []string{"location"})

	pressure = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "pressure",
		Help:      "Pressure in millibars",
	}, []string{"location"})

	precipitation = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "precipitation",
		Help:      "Precipitation amount in millimeters",
	}, []string{"location"})

	humidity = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "humidity",
		Help:      "Humidity as percentage",
	}, []string{"location"})
)
