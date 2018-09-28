package tracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

type Tags map[string]string

func SetUpTrace(ctx context.Context, opName string) func() {
	span, _ := opentracing.StartSpanFromContext(ctx, opName)
	return func() {
		span.Finish()
	}
}

func SetUpTraceWithTags(ctx context.Context, opName string, tags Tags) (deferFunc func()) {
	span, _ := opentracing.StartSpanFromContext(ctx, opName)
	for k, v := range tags {
		span.SetTag(k, v)
	}
	return func() {
		span.Finish()
	}
}

func ContextToString(ctx context.Context) string {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return ""
	}
	j, ok := span.Context().(jaeger.SpanContext)
	if !ok {
		return ""
	}
	return j.String()
}
