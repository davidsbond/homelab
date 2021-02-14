// Package scrape contains logic for scraping the health check endpoints of kubernetes pods that have specified
// the required annotations.
package scrape

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/opentracing/opentracing-go"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"pkg.dsb.dev/closers"
	"pkg.dsb.dev/health"
	"pkg.dsb.dev/tracing"
)

type (
	// The Scraper type is responsible for determining pods to scrape the health check endpoints of based on their
	// annotations.
	Scraper struct {
		k8s  *kubernetes.Clientset
		http *http.Client
		mux  *sync.RWMutex
		pods map[string]*PodHealth
	}

	// PodHealth represents the current health status of a given pod and its components.
	PodHealth struct {
		PodName     string       `json:"podName"`
		Namespace   string       `json:"namespace"`
		Name        string       `json:"name"`
		Description string       `json:"description"`
		Version     string       `json:"version"`
		Compiled    time.Time    `json:"compiled"`
		Healthy     bool         `json:"healthy"`
		Checks      []*Component `json:"checks"`
	}

	// The Component type represents a single aspect of a health check, this could be something like a database
	// connection, gRPC client etc.
	Component struct {
		Name    string `json:"name"`
		Healthy bool   `json:"healthy"`
		Message string `json:"message,omitempty"`
	}
)

// NewScraper returns a new instance of the Scraper type that will look for appropriately annotated pods in the cluster
// configured with the provided kubernetes.Clientset.
func NewScraper(k8s *kubernetes.Clientset) *Scraper {
	sc := &Scraper{
		k8s:  k8s,
		http: &http.Client{Timeout: time.Minute},
		mux:  &sync.RWMutex{},
		pods: make(map[string]*PodHealth),
	}

	health.AddCheck("k8s-scraper", sc.Ping)
	return sc
}

// Ping asserts that the connection with the k8s cluster works as expected by attempting to list pods in the
// default namespace.
func (s *Scraper) Ping() error {
	const ns = "default"
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	_, err := s.k8s.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{})
	return err
}

const (
	enabledAnnotation = "health.dsb.dev/scrape"
	portAnnotation    = "health.dsb.dev/port"
	pathAnnotation    = "health.dsb.dev/path"
	schemeAnnotation  = "health.dsb.dev/scheme"
	defaultPort       = "8081"
	defaultPath       = "/__/health"
	defaultScheme     = "http"
)

// Scrape pods with the appropriate annotations, returning a slice of PodHealth instances representing the current
// health check state of all annotated pods and their individual components.
func (s *Scraper) Scrape(ctx context.Context) ([]*PodHealth, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "health-scrape")
	defer span.Finish()

	s.mux.Lock()
	s.pods = make(map[string]*PodHealth)
	s.mux.Unlock()

	namespaces, err := s.getNamespaceNames(ctx)
	if err != nil {
		return nil, err
	}

	grp, ctx := errgroup.WithContext(ctx)
	for _, namespace := range namespaces {
		func(namespace string) {
			grp.Go(func() error {
				return s.scrapeNamespace(ctx, namespace)
			})
		}(namespace)
	}

	if err = grp.Wait(); err != nil {
		return nil, err
	}

	return s.results(), nil
}

func (s *Scraper) results() []*PodHealth {
	s.mux.RLock()
	defer s.mux.RUnlock()
	apps := make([]*PodHealth, len(s.pods))
	var i int
	for _, app := range s.pods {
		apps[i] = app
		i++
	}
	return apps
}

func (s *Scraper) getNamespaceNames(ctx context.Context) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "list-namespaces")
	defer span.Finish()

	list, err := s.k8s.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	names := make([]string, len(list.Items))
	for i, namespace := range list.Items {
		names[i] = namespace.Name
	}

	span.SetTag("count", len(names))
	return names, nil
}

func (s *Scraper) scrapeNamespace(ctx context.Context, namespace string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "scrape-namespace")
	defer span.Finish()
	span.SetTag("namespace", namespace)

	pods, err := s.listPods(ctx, namespace)
	if err != nil {
		return err
	}

	grp, ctx := errgroup.WithContext(ctx)
	for _, pod := range pods {
		if pod.Annotations[enabledAnnotation] != "true" {
			continue
		}
		if pod.Status.Phase != corev1.PodRunning {
			continue
		}

		func(pod corev1.Pod) {
			grp.Go(func() error {
				return s.scrapePod(ctx, pod)
			})
		}(pod)
	}

	return grp.Wait()
}

func (s *Scraper) listPods(ctx context.Context, namespace string) ([]corev1.Pod, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "list-pods")
	defer span.Finish()
	span.SetTag("namespace", namespace)

	pods, err := s.k8s.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	span.SetTag("count", len(pods.Items))
	return pods.Items, nil
}

func (s *Scraper) scrapePod(ctx context.Context, pod corev1.Pod) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "scrape-pod")
	defer span.Finish()

	span.SetTag("namespace", pod.Namespace)
	span.SetTag("name", pod.Name)

	endpoint := getPodEndpoint(pod)
	span.SetTag("endpoint", endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return tracing.WithError(span, err)
	}

	res, err := s.http.Do(req)
	if err != nil {
		err = fmt.Errorf("failed to request %s: %w", pod.Name, err)
		return tracing.WithError(span, err)
	}
	defer closers.Close(res.Body)

	var app PodHealth
	if err = json.NewDecoder(res.Body).Decode(&app); err != nil {
		err = fmt.Errorf("failed to decode health for %s: %w", pod.Name, err)
		return tracing.WithError(span, err)
	}

	app.PodName = pod.Name
	app.Namespace = pod.Namespace

	s.mux.Lock()
	defer s.mux.Unlock()
	s.pods[pod.Name] = &app
	return nil
}

func getPodEndpoint(pod corev1.Pod) string {
	ip := pod.Status.PodIP
	port := pod.Annotations[portAnnotation]
	if port == "" {
		port = defaultPort
	}
	path := pod.Annotations[pathAnnotation]
	if path == "" {
		path = defaultPath
	}
	scheme := pod.Annotations[schemeAnnotation]
	if scheme == "" {
		scheme = defaultScheme
	}

	return fmt.Sprintf("%s://%s:%s%s", scheme, ip, port, path)
}
