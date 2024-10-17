package otel

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func StartTracerSpanSingle(ctx context.Context, tracer trace.Tracer, name string, kind trace.SpanKind) (context.Context, trace.Span) {
	return tracer.Start(ctx, name, trace.WithSpanKind(kind))
}

func StartServerTracerSpan(ctx context.Context, tracer string, span string) (context.Context, trace.Span) {
	return StartTracerSpanSingle(ctx, otel.Tracer(tracer), span, trace.SpanKindServer)
}

func StartClientTracerSpan(ctx context.Context, tracer string, span string) (context.Context, trace.Span) {
	return StartTracerSpanSingle(ctx, otel.Tracer(tracer), span, trace.SpanKindClient)
}

func StartInternalTracerSpan(ctx context.Context, tracer string, span string) (context.Context, trace.Span) {
	return StartTracerSpanSingle(ctx, otel.Tracer(tracer), span, trace.SpanKindInternal)
}

func StartProducerTracerSpan(ctx context.Context, tracer string, span string) (context.Context, trace.Span) {
	return StartTracerSpanSingle(ctx, otel.Tracer(tracer), span, trace.SpanKindProducer)
}

func StartConsumerTracerSpan(ctx context.Context, tracer string, span string) (context.Context, trace.Span) {
	return StartTracerSpanSingle(ctx, otel.Tracer(tracer), span, trace.SpanKindConsumer)
}
