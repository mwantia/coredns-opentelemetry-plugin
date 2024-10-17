package opentelemetry

import (
	"context"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"

	otelsetup "github.com/mwantia/coredns-opentelemetry-plugin/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (p OtelPlugin) ServeDNS(ctx context.Context, writer dns.ResponseWriter, msg *dns.Msg) (int, error) {
	req := request.Request{W: writer, Req: msg}

	// Don't use the plugin name, since this will act as root
	ctx, span := otelsetup.StartDnsServerTracerSpan(ctx, req,
		"github.com/mwantia/coredns-opentelemetry-plugin", "ServeDNS")
	defer span.End()

	rw := dnstest.NewRecorder(writer)
	rcode, err := plugin.NextOrFailure(p.Name(), p.Next, ctx, rw, msg)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	span.SetAttributes(
		attribute.String("dns.rcode", dns.RcodeToString[rcode]),
	)

	return rcode, err
}
