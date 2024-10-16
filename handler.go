package otel

import (
	"context"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type OtelHandler struct {
	Tracer trace.Tracer
	Next   plugin.Handler
}

func (h *OtelHandler) Name() string {
	return "tracing:" + h.Next.Name()
}

func (h *OtelHandler) ServeDNS(ctx context.Context, writer dns.ResponseWriter, msg *dns.Msg) (int, error) {
	name := h.Next.Name()

	ctx, span := h.Tracer.Start(ctx, name)
	defer span.End()

	status, err := h.Next.ServeDNS(ctx, writer, msg)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	span.SetAttributes(
		attribute.Int("dns.status", status),
	)

	return status, err
}
