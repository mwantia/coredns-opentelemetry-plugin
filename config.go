package otel

import (
	"os"

	"github.com/coredns/caddy"
	"github.com/rschone/corefile2struct/pkg/corefile"
)

type OtelConfig struct {
	Endpoint    string `cf:"endpoint" check:"nonempty"`
	ServiceName string `cf:"servicename" default:"coredns"`
	Hostname    string `cf:"hostname"`
}

func CreatePlugin(c *caddy.Controller) (*OtelPlugin, error) {
	cfg, err := ParseConfig(c)
	if err != nil {
		return nil, err
	}

	return &OtelPlugin{
		Cfg: *cfg,
	}, nil
}

func ParseConfig(c *caddy.Controller) (*OtelConfig, error) {
	var cfg OtelConfig
	if err := corefile.Parse(c, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (p OtelPlugin) GetHostname() string {
	if len(p.Cfg.Hostname) <= 0 {
		hostname, _ := os.Hostname()
		return hostname
	}

	return p.Cfg.Hostname
}
