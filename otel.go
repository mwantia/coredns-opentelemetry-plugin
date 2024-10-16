package otel

import (
	"context"
	"fmt"
	"time"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

func (p *OtelPlugin) OnStartup() error {
	exporter, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpoint(p.Cfg.Endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("failed to create OTLP exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter,
			sdktrace.WithBatchTimeout(5*time.Second),
			sdktrace.WithMaxExportBatchSize(10),
		),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(p.Cfg.ServiceName),
			semconv.TelemetrySDKNameKey.String("opentelemetry"),
			semconv.TelemetrySDKLanguageKey.String("go"),
			attribute.String("hostname", p.GetHostname()),
		)),
	)

	otel.SetTracerProvider(tp)
	p.Tracer = otel.Tracer("coredns.otel")

	return nil
}

func (p OtelPlugin) ServeDNS(ctx context.Context, writer dns.ResponseWriter, msg *dns.Msg) (int, error) {
	req := request.Request{W: writer, Req: msg}

	// Don't use the plugin name, since this will act as root
	ctx, span := p.Tracer.Start(ctx, "ServeDNS",
		trace.WithAttributes(
			attribute.String("dns.name", req.Name()),
			attribute.String("dns.type", req.Type()),
			attribute.String("dns.proto", req.Proto()),
			attribute.String("dns.remote", req.IP()),
		),
		trace.WithSpanKind(trace.SpanKindServer),
	)
	defer span.End()

	rw := dnstest.NewRecorder(writer)
	status, err := plugin.NextOrFailure(p.Name(), p.Next, ctx, rw, msg)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	span.SetAttributes(
		attribute.Int("dns.rcode", rw.Rcode),
		attribute.Int("dns.status", status),
	)

	return status, err
}
