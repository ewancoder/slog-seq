package slogseq

import (
	"context"

	"go.opentelemetry.io/otel/sdk/trace"
)

type LoggingSpanProcessor struct {
	Handler *SeqHandler // or a slog.Logger that wraps your SeqHandler
}

func (p *LoggingSpanProcessor) OnStart(ctx context.Context, s trace.ReadWriteSpan) {
	// no-op, or you can log the start if you like
}

func (p *LoggingSpanProcessor) OnEnd(s trace.ReadOnlySpan) {
	// Called when the span ends
	events := s.Events()
	for _, e := range events {
		// e.Name, e.Time, e.Attributes, e.DroppedAttributeCount
		// Convert these into log or CLEF events
		// For example:
		p.logOtelEventAsCLEF(s, e)
	}
}

func (p *LoggingSpanProcessor) ForceFlush(ctx context.Context) error {
	// flush logs if needed
	return nil
}

func (p *LoggingSpanProcessor) Shutdown(ctx context.Context) error {
	// gracefully close
	return nil
}

func (p *LoggingSpanProcessor) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	// Export spans if needed
	for _, s := range spans {
		for _, e := range s.Events() {
			p.logOtelEventAsCLEF(s, e)
		}
	}
	return nil
}

// logOtelEventAsCLEF converts an OTEL event into a CLEF log event
func (p *LoggingSpanProcessor) logOtelEventAsCLEF(span trace.ReadOnlySpan, e trace.Event) {
	sc := span.SpanContext()
	if !sc.IsValid() {
		return
	}

	event := &CLEFEvent{
		Timestamp:          e.Time,
		Message:            e.Name,
		TraceID:            sc.TraceID().String(),
		SpanID:             sc.SpanID().String(),
		SpanStart:          span.StartTime(),
		ResourceAttributes: map[string]interface{}{"service": map[string]interface{}{"name": span.Name()}},
		Properties:         make(map[string]interface{}),
	}

	if parent := span.Parent(); parent.IsValid() {
		event.ParentSpanID = parent.SpanID().String()
	}

	for _, attr := range e.Attributes {
		k := string(attr.Key)
		v := attr.Value.AsInterface()
		event.Properties[k] = v
		if k == "exception.message" {
			event.Level = "ERROR"
			event.Message = v.(string)
		}
	}

	p.Handler.HandleCLEFEvent(*event)
}
