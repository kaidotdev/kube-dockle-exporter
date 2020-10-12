package collector

import (
	"context"
	"encoding/json"
	"kube-dockle-exporter/pkg/client"
	"sync"
	"time"

	"golang.org/x/xerrors"

	"github.com/prometheus/client_golang/prometheus"
	v1 "k8s.io/api/core/v1"
)

const (
	namespace = "dockle"
)

type DockleCollector struct {
	Logger           ILogger
	KubernetesClient IKubernetesClient
	DockleClient     IDockleClient
	concurrency      int64
	vulnerabilities  *prometheus.GaugeVec
}

func NewDockleCollector(
	logger ILogger,
	kubernetesClient IKubernetesClient,
	dockleClient IDockleClient,
	concurrency int64,
) *DockleCollector {
	return &DockleCollector{
		Logger:           logger,
		KubernetesClient: kubernetesClient,
		DockleClient:     dockleClient,
		concurrency:      concurrency,
		vulnerabilities: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "cis_benchmarks_total",
			Help:      "CIS benchmarks executed by dockle",
		}, []string{"image", "code", "level"}),
	}
}

func uniqueContainerImages(containers []v1.Container) []string {
	keys := make(map[string]bool)
	var images []string
	for _, container := range containers {
		image := container.Image
		if _, value := keys[image]; !value {
			keys[image] = true
			images = append(images, image)
		}
	}
	return images
}

func (c *DockleCollector) Scan(ctx context.Context) error {
	containers, err := c.KubernetesClient.Containers()
	if err != nil {
		return xerrors.Errorf("failed to get containers: %w", err)
	}

	semaphore := make(chan struct{}, c.concurrency)
	defer close(semaphore)

	wg := sync.WaitGroup{}
	mutex := &sync.Mutex{}

	var dockleResponses []client.DockleResponse
	for _, image := range uniqueContainerImages(containers) {
		wg.Add(1)
		go func(image string) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() {
				<-semaphore
			}()
			out, err := c.DockleClient.Do(ctx, image)
			if err != nil {
				c.Logger.Errorf("Failed to execute CIS benchmark at %s: %s\n", image, err.Error())
				return
			}

			var response client.DockleResponse
			if err := json.Unmarshal(out, &response); err != nil {
				c.Logger.Errorf("Failed to parse dockle response at %s: %s\n", image, err.Error())
				return
			}
			response.Target = image
			func() {
				mutex.Lock()
				defer mutex.Unlock()
				dockleResponses = append(dockleResponses, response)
			}()
		}(image)
	}
	wg.Wait()

	c.vulnerabilities.Reset()
	for _, dockleResponse := range dockleResponses {
		for _, detail := range dockleResponse.Details {
			labels := []string{
				dockleResponse.ExtractImage(),
				detail.Code,
				detail.Level,
			}
			c.vulnerabilities.WithLabelValues(labels...).Set(1)
		}
	}

	return nil
}

func (c *DockleCollector) StartLoop(ctx context.Context, interval time.Duration) {
	go func(ctx context.Context) {
		t := time.NewTicker(interval)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				if err := c.Scan(ctx); err != nil {
					c.Logger.Errorf("Failed to scan: %s\n", err.Error())
				}
			case <-ctx.Done():
				return
			}
		}
	}(ctx)
}

func (c *DockleCollector) collectors() []prometheus.Collector {
	return []prometheus.Collector{
		c.vulnerabilities,
	}
}

func (c *DockleCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, collector := range c.collectors() {
		collector.Describe(ch)
	}
}

func (c *DockleCollector) Collect(ch chan<- prometheus.Metric) {
	for _, collector := range c.collectors() {
		collector.Collect(ch)
	}
}
