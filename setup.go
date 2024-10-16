package otel

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/mwantia/coredns-otel-plugin/logging"
	"github.com/mwantia/coredns-otel-plugin/metrics"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const PluginName string = "opentelemetry"

type OtelPlugin struct {
	Next     plugin.Handler
	Cfg      OtelConfig
	Provider *sdktrace.TracerProvider
	Tracer   trace.Tracer
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
		logging.Log.Errorf("%v", err)
		return plugin.Error(PluginName, err)
	}

	c.OnStartup(p.OnStartup)
	c.OnShutdown(p.OnShutdown)

	if err := metrics.Register(); err != nil {
		return err
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		p.Next = next
		return p
	})

	return nil
}
