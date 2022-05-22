package otelexample

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

// StartExample uses an OpenTelemetry tracer to start a span. This is more or less straight
// from the documentation:
// https://opentelemetry.io/docs/instrumentation/go/manual/#creating-spans
//
// Since Tracer.Start is a call through an interface, all the arguments escape to the heap.
func StartExample(ctx context.Context) {
	// otelexample.go:20:28: inlining call to attribute.String
	// otelexample.go:20:28: inlining call to attribute.Key.String
	// otelexample.go:20:28: inlining call to attribute.StringValue
	spanKV := attribute.String("span_key", "span_value")
	// otelexample.go:25:75: inlining call to trace.WithAttributes
	// otelexample.go:25:27: ... argument escapes to heap
	// otelexample.go:25:75: ... argument escapes to heap
	// otelexample.go:25:75: trace.SpanStartEventOption(trace.attributeOption(trace.attributes)) escapes to heap
	ctx, span := tracer.Start(ctx, "operation-with-arg", trace.WithAttributes(spanKV))
	defer span.End()

	fmt.Printf("do something ctx=%v\n", ctx)
}
