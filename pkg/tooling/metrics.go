package tooling

import "github.com/prometheus/client_golang/prometheus"

func NewRegistry() *prometheus.Registry {
	registry := prometheus.NewRegistry()
	return registry
}
