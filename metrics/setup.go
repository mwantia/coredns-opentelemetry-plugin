package metrics

import "github.com/prometheus/client_golang/prometheus"

func Register() error {
	prometheus.MustRegister(metricsQueryRequestsTotal)
	return nil
}
