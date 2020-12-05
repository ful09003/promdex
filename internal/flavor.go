package internal

import (
	"github.com/prometheus/client_golang/api/prometheus/v1"

	"time"
)

//PromdexFlavor holds data pertaining to Promdex 'flavor' text
type PromdexFlavor struct {
	CtxString string `json:"context"` //String representing the context (flavor) for a particular Prom metric
}

//NewFlavor takes a string containing metric context and returns a PromdexFlavor
func NewFlavor(c string) PromdexFlavor {
	var f PromdexFlavor
	f.CtxString = c

	return f
}

//PromdexEnhancedMetric is a union of a Prometheus metric metadata and Promdex flavor
type PromdexEnhancedMetric struct {
	PromMetric v1.MetricMetadata `json:"prometheus_metric_data"`
	PromdexMetric []PromdexFlavor `json:"promdex_metric_data"`
	Job string `json:"prometheus_job_name"`
	Created time.Time `json:"promdex_metric_created_time"`
	Updated time.Time `json:"promdex_metric_updated_time"`
}