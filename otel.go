package opentelemetry

import (
	"context"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func (p OtelPlugin) ServeDNS(ctx context.Context, writer dns.ResponseWriter, msg *dns.Msg) (int, error) {
	req := request.Request{W: writer, Req: msg}

	// Don't use the plugin name, since this will act as root
	tracer := otel.Tracer("github.com/mwantia/coredns-opentelemetry-plugin")
	ctx, span := tracer.Start(ctx, "ServeDNS",
		trace.WithAttributes(
			attribute.String("dns.fqdn", dns.Fqdn(req.Name())),
			attribute.String("dns.type", req.Type()),
			attribute.String("dns.proto", req.Proto()),
			attribute.String("dns.remote", req.IP()),
		),
		trace.WithSpanKind(trace.SpanKindServer),
	)
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
