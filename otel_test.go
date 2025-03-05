package slogseq

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func TestOnEnd_WithException(t *testing.T) {
	handler := &SeqHandler{noFlush: true, workerCount: 1}
	handler.start()
	processor := &LoggingSpanProcessor{Handler: handler}

	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(processor))
	// Ensure spans are processed before the test exits.
	defer func() { _ = tp.Shutdown(context.Background()) }()

	// Obtain a tracer from the provider.
	tracer := tp.Tracer("test-tracer")

	// Start a span, add an event with exception attributes, then end the span.
	ctx := context.Background()
	_, span := tracer.Start(ctx, "testSpan", trace.WithSpanKind(trace.SpanKindServer))
	span.AddEvent("originalEventName", trace.WithAttributes(
		attribute.String("exception.message", "error occurred"),
		attribute.Int("code", 500),
	))
	span.End()

	var evt CLEFEvent

	select {
	case evt = <-handler.workers[0].eventsCh:
	case <-time.After(1000 * time.Millisecond):
		t.Fatal("timed out waiting for event")
	}

	// Check that the exception message overwrote the event's original name.
	if evt.Message != "error occurred" {
		t.Errorf("expected message 'error occurred', got %s", evt.Message)
	}
	// Check that the level was set to error.
	if evt.Level != CLEFLevelError.String() {
		t.Errorf("expected level %s, got %s", CLEFLevelError.String(), evt.Level)
	}
	// Check that additional properties (like code) are present.
	if code, ok := evt.Properties["code"]; !ok {
		t.Errorf("expected property 'code' to be set")
	} else if code.(int64) != 500 {
		t.Errorf("expected code 500, got %v", code)
	}
}
