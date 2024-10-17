package opentelemetry

import (
	"os"
	"strconv"
	"time"

	"github.com/coredns/caddy"
	"github.com/rschone/corefile2struct/pkg/corefile"
)

type OtelConfig struct {
	Endpoint     string `cf:"endpoint" check:"nonempty"`
	ServiceName  string `cf:"servicename" default:"coredns"`
	Hostname     string `cf:"hostname"`
	BatchTimeout string `cf:"batchtimeout" check:"nonempty" default:"5s"`
	BatchSize    string `cf:"batchsize" check:"nonempty" default:"10"`
	SamplingRate string `cf:"samplingrate" check:"nonempty" default:"5"`
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

func (p OtelPlugin) GetBatchTimeout() time.Duration {
	timeout, err := time.ParseDuration(p.Cfg.BatchTimeout)
	if err != nil {
		// Default is 5 seconds
		return 5 * time.Second
	}

	return timeout
}

func (p OtelPlugin) GetBatchSize() int {
	size, err := strconv.Atoi(p.Cfg.BatchSize)
	if err != nil {
		// Default is 10
		return 10
	}

	return size
}

func (p OtelPlugin) GetSamplingRateFraction() float64 {
	rate, err := strconv.Atoi(p.Cfg.SamplingRate)
	if err != nil {
		// Default is 5
		return float64(5)
	}

	return float64(rate)
}
