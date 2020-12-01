package internal

import (
	"context"
	"net/url"
	"sync"
	"time"

	api "github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	log "github.com/sirupsen/logrus"
)

//TargetDiscoverer is our set of behaviors for discovering Prometheus data
type TargetDiscoverer interface {
	QueryTargetMetadata() ([]v1.MetricMetadata, error)
	GetTarget() string
}

//PromdexTarget represents a Prometheus endpoint which we scrape metrics metadata from
type PromdexTarget struct {
	Target         string
	LastDiscovered time.Time
	client         api.Client

	sync.Mutex
}

//NewPromdexTarget takes a target URL representing a Prometheus endpoint and returns a PromdexTarget
func NewPromdexTarget(t string) (*PromdexTarget, error) {
	var r PromdexTarget

	if err := isTargetValid(t); err != nil {
		return &r, err
	}

	r.Target = t
	r.LastDiscovered = time.Now().UTC()
	c, e := api.NewClient(api.Config{
		Address: t,
	})
	if e != nil {
		return &r, e
	}
	r.client = c

	return &r, nil
}

func isTargetValid(t string) error {
	_, e := url.ParseRequestURI(t)

	if e != nil {
		return e
	}

	return nil
}

//QueryTargetMetadata queries the Promdex target's metadata (akin to a 'live' lookup), returning the raw metadata and/or any errors encountered
func (p *PromdexTarget) QueryTargetMetadata() ([]v1.MetricMetadata, error) {
	p.Lock()
	defer p.Unlock()

	api := v1.NewAPI(p.client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	metadata, err := api.TargetsMetadata(ctx, "", "", "")
	if err != nil {
		return nil, err
	}
	return metadata, err
}

//GetTarget returns the PromdexTarget's target string
func (p *PromdexTarget) GetTarget() string {
	return p.Target
}

//PromdexDiscoverStorer is a struct implementing the TargetDiscoverer and Metastorer interfaces
type PromdexDiscoverStorer struct {
	TargetDiscoverer
	Metastorer
}

// LoadIn causes the PromdexDiscoverStorer to load metrics metadata from its configured target and persist those to its Metastorer
func (p PromdexDiscoverStorer) LoadIn() (int, error) {
	log.WithFields(log.Fields{
		"targetURL": p.GetTarget(),
	}).Debug("beginning metrics loading")
	metrics, err := p.TargetDiscoverer.QueryTargetMetadata()

	if err != nil {
		log.WithFields(log.Fields{
			"targetURL": p.GetTarget(),
		}).Errorf("failed loading metrics, original: %s", err)
		return 0, err
	}

	for _, m := range metrics {
		if err := p.Metastorer.StoreExporterMeta(m.Target["job"], m.Metric, m.Help); err != nil {
			log.WithFields(log.Fields{
				"metricName": m.Metric,
			}).Errorf("failed writing metric, original: %s", err)
		}
	}

	return len(metrics), nil
}
