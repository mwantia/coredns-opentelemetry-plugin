package opentelemetry

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
	"github.com/mwantia/coredns-opentelemetry-plugin/logging"
	otelsetup "github.com/mwantia/coredns-opentelemetry-plugin/otel"
)

func TestOtel(tst *testing.T) {
	OverwriteStdOut()

	c := caddy.NewTestController("dns", `
		opentelemetry {
			endpoint jaeger:4318
			hostname localhost
			samplingrate 1
		}
	`)
	cfg, err := ParseConfig(c)
	if err != nil {
		tst.Errorf("Unable to parse config: %v", err)
	}

	p := OtelPlugin{
		Cfg: *cfg,
	}

	ctx := context.TODO()
	shutdown, err := otelsetup.SetupOpentelemetry(ctx, otelsetup.OpenTelemtryConfig{
		Endpoint:     p.Cfg.Endpoint,
		ServiceName:  p.Cfg.ServiceName,
		Hostname:     p.GetHostname(),
		BatchTimeout: p.GetBatchTimeout(),
		BatchSize:    p.GetBatchSize(),
		SamplingRate: p.GetSamplingRateFraction(),
	})

	RunTests(ctx, tst, p)

	if shutdown != nil {
		// time.Sleep(10 * time.Second)
		err := shutdown(ctx)
		if err != nil {
			tst.Errorf("Unable to gracefully shutdown plugin: %v", err)
		}
	}
}

func RunTests(ctx context.Context, tst *testing.T, p OtelPlugin) {
	tst.Run("", func(t *testing.T) {
		req := new(dns.Msg)
		req.SetQuestion("google.de", dns.TypeA)
		rec := dnstest.NewRecorder(&test.ResponseWriter{})

		code, err := p.ServeDNS(ctx, rec, req)

		logging.Log.Infof("Code:  %v", code)
		logging.Log.Infof("Error: %v", err)
	})
}

func OverwriteStdOut() error {
	tempFile, err := os.CreateTemp("", "coredns-consulkv-test-log")
	if err != nil {
		return err
	}

	defer os.Remove(tempFile.Name())

	orig := logging.Log
	logging.Log = clog.NewWithPlugin("consulkv")
	log.SetOutput(os.Stdout)

	defer func() {
		logging.Log = orig
	}()

	clog.D.Set()
	return nil
}
