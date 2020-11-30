package internal

import (
	"net/url"
	"time"
	"sync"
	"context"
	api "github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

//PromdexTarget represents a Prometheus endpoint which we scrape metrics metadata from
type PromdexTarget struct {
	Target         string
	LastDiscovered time.Time
	client api.Client

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
func (p *PromdexTarget) QueryTargetMetadata() (map[string][]v1.Metadata, error) {
	p.Lock()
	defer p.Unlock()

	api := v1.NewAPI(p.client)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	metadata, err := api.Metadata(ctx, "", "")
	if err != nil {
		return nil, err
	}
	return metadata, err
}