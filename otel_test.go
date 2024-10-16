package otel

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
	"github.com/mwantia/coredns-otel-plugin/logging"
)

func TestOtel(tst *testing.T) {
	OverwriteStdOut()

	ctx := context.TODO()
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

	if err := p.OnStartup(); err != nil {
		tst.Errorf("Unable to start plugin: %v", err)
	}

	RunTests(ctx, tst, p)

	if err := p.OnShutdown(); err != nil {
		tst.Errorf("Unable to gracefully shutdown plugin: %v", err)
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
