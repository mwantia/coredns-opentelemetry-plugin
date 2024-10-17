package opentelemetry

import (
	"context"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/mwantia/coredns-opentelemetry-plugin/metrics"
	"github.com/mwantia/coredns-opentelemetry-plugin/otel"
)

const PluginName string = "opentelemetry"

type OtelPlugin struct {
	Next plugin.Handler
	Cfg  OtelConfig
}

func init() {
	plugin.Register(PluginName, setup)
}

func (p OtelPlugin) Name() string {
	return PluginName
}

func setup(c *caddy.Controller) error {
	p, err := CreatePlugin(c)
	if err != nil {
		return plugin.Error(PluginName, err)
	}

	var shutdown func(ctx context.Context) error

	c.OnStartup(func() error {
		if err := metrics.Register(); err != nil {
			return err
		}

		ctx := context.Background()
		shutdown, err = otel.SetupOpentelemetry(ctx, otel.OpenTelemtryConfig{
			Endpoint:     p.Cfg.Endpoint,
			ServiceName:  p.Cfg.ServiceName,
			Hostname:     p.GetHostname(),
			BatchTimeout: p.GetBatchTimeout(),
			BatchSize:    p.GetBatchSize(),
			SamplingRate: p.GetSamplingRateFraction(),
		})

		return err
	})

	c.OnShutdown(func() error {
		if shutdown != nil {
			ctx := context.Background()
			return shutdown(ctx)
		}

		return nil
	})

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		p.Next = next
		return p
	})

	return nil
}
