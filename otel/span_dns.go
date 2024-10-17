package otel

import (
	"context"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func StartTracerSpanDns(ctx context.Context, tracer trace.Tracer, name string, req request.Request, kind trace.SpanKind) (context.Context, trace.Span) {
	return tracer.Start(ctx, name,
		trace.WithAttributes(
			attribute.String("dns.fqdn", dns.Fqdn(req.Name())),
			attribute.String("dns.type", req.Type()),
			attribute.String("dns.proto", req.Proto()),
			attribute.String("dns.remote", req.IP()),
		),
		trace.WithSpanKind(kind),
	)
}

func StartDnsServerTracerSpan(ctx context.Context, req request.Request, tracer, span string) (context.Context, trace.Span) {
	return StartTracerSpanDns(ctx, otel.Tracer(tracer), span, req, trace.SpanKindServer)
}

func StartDnsClientTracerSpan(ctx context.Context, req request.Request, tracer, span string) (context.Context, trace.Span) {
	return StartTracerSpanDns(ctx, otel.Tracer(tracer), span, req, trace.SpanKindClient)
}

func StartDnsInternalTracerSpan(ctx context.Context, req request.Request, tracer, span string) (context.Context, trace.Span) {
	return StartTracerSpanDns(ctx, otel.Tracer(tracer), span, req, trace.SpanKindInternal)
}

func StartDnsProducerTracerSpan(ctx context.Context, req request.Request, tracer, span string) (context.Context, trace.Span) {
	return StartTracerSpanDns(ctx, otel.Tracer(tracer), span, req, trace.SpanKindProducer)
}

func StartDnsConsumerTracerSpan(ctx context.Context, req request.Request, tracer, span string) (context.Context, trace.Span) {
	return StartTracerSpanDns(ctx, otel.Tracer(tracer), span, req, trace.SpanKindConsumer)
}

func StartDnsTestTracerSpan(ctx context.Context, dnsName, dnsType, tracerName string) (context.Context, trace.Span) {
	return otel.Tracer(tracerName).Start(ctx, "TestDNS",
		trace.WithAttributes(
			attribute.String("dns.fqdn", dns.Fqdn(dnsName)),
			attribute.String("dns.type", dnsType),
			attribute.String("dns.proto", "udp"),
			attribute.String("dns.remote", "localhost"),
		),
		trace.WithSpanKind(trace.SpanKindServer),
	)
}
