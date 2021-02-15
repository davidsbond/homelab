// Package directory contains logic for serving the homelab index page.
package directory

import (
	"html/template"
	"net/http"
	"os"
	"sort"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
	"pkg.dsb.dev/closers"
	"pkg.dsb.dev/transport"

	"github.com/davidsbond/homelab/internal/directory/assets"
)

//go:generate go-bindata -pkg assets -prefix assets -nometadata -ignore bindata -o ./assets/bindata.go ./assets

type (
	// The HTTP type is responsible for handling HTTP requests for the index page.
	HTTP struct {
		transport.HTTP

		config *Config
	}

	// The Config type represents the YAML configuration file for links to display on the index page.
	Config struct {
		Items []Item `yaml:"items"`
	}

	// The Item type represents a single item to display on the index page.
	Item struct {
		Name        string `yaml:"name"`
		URL         string `yaml:"url"`
		Description string `yaml:"description"`
	}
)

// LoadConfig attempts to load the index page configuration located at the provided path. Returns an error if it
// does not exist or is an invalid YAML configuration.
func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer closers.Close(f)

	var config Config
	return &config, yaml.NewDecoder(f).Decode(&config)
}

// NewHTTP returns a new instance of the HTTP type that will render the index page with the provided configuration
// items. Items are sorted alphabetically by name.
func NewHTTP(config *Config) *HTTP {
	sort.Slice(config.Items, func(i, j int) bool {
		return config.Items[i].Name < config.Items[j].Name
	})

	return &HTTP{config: config}
}

// Register HTTP handlers on the provided mux.Router.
func (t *HTTP) Register(r *mux.Router) {
	r.HandleFunc("/", t.Get).Methods(http.MethodGet)
}

// Get handles an inbound HTTP GET request that renders the index page template with the provided configuration
// and writes the resulting HTML to the client.
func (t *HTTP) Get(w http.ResponseWriter, r *http.Request) {
	const file = "index.gohtml"

	data, err := assets.Asset(file)
	if err != nil {
		t.Error(r.Context(), w, http.StatusInternalServerError, err.Error())
		return
	}

	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
	}

	tpl, err := template.New("").Funcs(funcMap).Parse(string(data))
	if err != nil {
		t.Error(r.Context(), w, http.StatusInternalServerError, err.Error())
		return
	}

	if err = tpl.Execute(w, t.config); err != nil {
		t.Error(r.Context(), w, http.StatusInternalServerError, err.Error())
		return
	}
}
