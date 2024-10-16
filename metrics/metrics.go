package metrics

import (
	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
	"github.com/prometheus/client_golang/prometheus"
)

var metricsQueryRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: plugin.Namespace,
	Subsystem: "otel",
	Name:      "query_requests_total",
	Help:      "Count the amount of queries received as request by the plugin.",
}, []string{"server", "type"})

func MetricsQueryRequests(server string, qtype uint16) {
	t := dns.TypeToString[qtype]
	metricsQueryRequestsTotal.WithLabelValues(server, t).Inc()
}
