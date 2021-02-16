// Package health contains an HTTP server that renders a basic HTML UI displaying the results of health
// checks of kubernetes pods that have special annotations.
package health

import (
	"context"
	"net/http"
	"sort"
	"text/template"

	"github.com/gorilla/mux"
	"pkg.dsb.dev/transport"

	"github.com/davidsbond/homelab/internal/health/assets"
	"github.com/davidsbond/homelab/internal/health/scrape"
)

//go:generate go-bindata -pkg assets -prefix assets -nometadata -ignore bindata -o ./assets/bindata.go ./assets

type (
	// The HTTP type is responsible for handling inbound HTTP requests to display the health check user-interface.
	HTTP struct {
		transport.HTTP

		scraper Scraper
	}

	// The Scraper interface describes types that scrape kubernetes pods for their current health.
	Scraper interface {
		Scrape(ctx context.Context) ([]*scrape.PodHealth, error)
	}
)

// NewHTTP returns a new instance of the HTTP type that will scrape pod health using the provided Scraper
// implementation.
func NewHTTP(scraper Scraper) *HTTP {
	return &HTTP{scraper: scraper}
}

// Register HTTP handlers on the provided mux.Router.
func (t *HTTP) Register(r *mux.Router) {
	r.HandleFunc("/", t.Get).Methods(http.MethodGet)
}

// Get handles an inbound HTTP GET request that performs a health scrape and renders the health check UI.
func (t *HTTP) Get(w http.ResponseWriter, r *http.Request) {
	const file = "index.gohtml"

	data, err := assets.Asset(file)
	if err != nil {
		t.Error(r.Context(), w, http.StatusInternalServerError, err.Error())
		return
	}

	tpl, err := template.New("").Parse(string(data))
	if err != nil {
		t.Error(r.Context(), w, http.StatusInternalServerError, err.Error())
		return
	}

	apps, err := t.scraper.Scrape(r.Context())
	if err != nil {
		t.Error(r.Context(), w, http.StatusInternalServerError, err.Error())
		return
	}

	sort.Slice(apps, func(i, j int) bool {
		return apps[i].Name < apps[j].Name
	})

	if err = tpl.Execute(w, apps); err != nil {
		t.Error(r.Context(), w, http.StatusInternalServerError, err.Error())
		return
	}
}
