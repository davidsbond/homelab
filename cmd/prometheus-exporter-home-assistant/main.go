package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/mux"
	"pkg.dsb.dev/app"
	"pkg.dsb.dev/closers"
	"pkg.dsb.dev/flag"
	"pkg.dsb.dev/server"
)

func main() {
	a := app.New(
		app.WithRunner(run),
		app.WithFlags(
			&flag.String{
				Name:        "home-assistant-url",
				Usage:       "Base URL of the home-assistant instance",
				EnvVar:      "HOME_ASSISTANT_URL",
				Required:    true,
				Destination: &homeAssistantURL,
			},
			&flag.String{
				Name:        "home-assistant-token",
				Usage:       "The long-lived access token to use to authenticate",
				EnvVar:      "HOME_ASSISTANT_TOKEN",
				Required:    true,
				Destination: &homeAssistantToken,
			},
		),
	)

	if err := a.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

var (
	homeAssistantURL   string
	homeAssistantToken string
)

func run(ctx context.Context) error {
	u, err := url.Parse(homeAssistantURL)
	if err != nil {
		return err
	}

	u.Path = "/api/prometheus"

	r := mux.NewRouter()
	r.HandleFunc("/", homeAssistantHandler(u))

	return server.ServeHTTP(ctx, server.WithHandler(r))
}

func homeAssistantHandler(u fmt.Stringer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := &http.Client{Timeout: time.Minute}

		req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, u.String(), nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req.Header.Set("Authorization", "Bearer "+homeAssistantToken)

		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer closers.Close(resp.Body)

		if _, err = io.Copy(w, resp.Body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(resp.StatusCode)
	}
}
