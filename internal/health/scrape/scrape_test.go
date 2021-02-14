package scrape_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"pkg.dsb.dev/environment"

	"github.com/davidsbond/homelab/internal/health/scrape"
)

func TestChecker_Check(t *testing.T) {
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	assert.NoError(t, err)
	assert.NotNil(t, config)

	client, err := k8s.NewForConfig(config)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	ctx := environment.NewContext()
	scraper := scrape.NewScraper(client)
	apps, err := scraper.Scrape(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, apps)
}
