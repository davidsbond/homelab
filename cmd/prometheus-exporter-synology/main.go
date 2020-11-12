package main

import (
	"context"
	"fmt"
	"os"

	"pkg.dsb.dev/app"
	"pkg.dsb.dev/flag"
	"pkg.dsb.dev/metrics"

	"github.com/davidsbond/homelab/internal/synology"
)

func main() {
	a := app.New(
		app.WithRunner(run),
		app.WithFlags(
			&flag.String{
				Name:        "synology-url",
				Usage:       "The URL of the Synology NAS",
				EnvVar:      "SYNOLOGY_URL",
				Required:    true,
				Destination: &synologyURL,
			},
			&flag.String{
				Name:        "synology-user",
				Usage:       "The username to use for authentication",
				EnvVar:      "SYNOLOGY_USER",
				Required:    true,
				Destination: &synologyUser,
			},
			&flag.String{
				Name:        "synology-password",
				Usage:       "The password to use for authentication",
				EnvVar:      "SYNOLOGY_PASSWORD",
				Required:    true,
				Destination: &synologyPass,
			},
		),
	)

	if err := a.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

var (
	synologyURL  string
	synologyUser string
	synologyPass string
)

func run(ctx context.Context) error {
	metrics.Register(diskTemp, diskSize, volumeSize, volumeUsed, uptime)

	cl, err := synology.New(synologyURL, synologyUser, synologyPass)
	if err != nil {
		return err
	}

	info, err := cl.SystemInfo(ctx)
	if err != nil {
		return err
	}

	for _, disk := range info.Disks {
		diskTemp.WithLabelValues(disk.Name).Set(disk.Temperature)
		diskSize.WithLabelValues(disk.Name).Set(disk.Size)
	}

	for _, volume := range info.Volumes {
		volumeSize.WithLabelValues(volume.Name).Set(volume.Size)
		volumeUsed.WithLabelValues(volume.Name).Set(volume.Used)
	}

	uptime.Set(info.Uptime)
	return metrics.Push()
}
