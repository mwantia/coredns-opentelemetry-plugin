package otel

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/mwantia/coredns-otel-plugin/logging"
	"github.com/mwantia/coredns-otel-plugin/metrics"
	"go.opentelemetry.io/otel/trace"
)

type OtelPlugin struct {
	Next   plugin.Handler
	Cfg    OtelConfig
	Tracer trace.Tracer
}

func init() {
	plugin.Register("otel", setup)
}

func (p OtelPlugin) Name() string {
	return "otel"
}

func setup(c *caddy.Controller) error {
	p, err := CreatePlugin(c)
	if err != nil {
		logging.Log.Errorf("%v", err)
		return plugin.Error("netboxdns", err)
	}

	c.OnStartup(func() error {
		if err := metrics.Register(); err != nil {
			return err
		}
		if err := p.OnStartup(); err != nil {
			return err
		}

		return nil
	})

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		p.Next = next
		return p
	})

	return nil
}
