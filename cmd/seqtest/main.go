package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"path"
	"time"

	slogseq "github.com/sokkalf/slog-seq" // import your library
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
	tr "go.opentelemetry.io/otel/trace"
)

var (
	seqURL = flag.String("url", "http://localhost:5341/ingest/clef", "Seq ingestion URL")
	apiKey = flag.String("key", "", "Seq API key")
)

func main() {
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		return
	}
	opts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "password" {
				a.Value = slog.StringValue("*****")
			}
			if a.Key == slog.SourceKey {
				// Replace the full path with just the file name
				s := a.Value.Any().(*slog.Source)
				s.File = path.Base(s.File)
			}
			return a
		},
	}

	seqLogger, handler := slogseq.NewSeqLogger(
		*seqURL,       // seqURL
		*apiKey,       // apiKey
		50,            // batchSize
		2*time.Second, // flushInterval
		opts,          // opts
	)
	defer handler.Close()

	slog.SetDefault(seqLogger.With("app", "slog-seq").With("env", "dev").With("version", "1.0.0"))

	slog.Info("Hello from slog-seq test command!",
		"env", "dev",
		"version", "1.0.0")

	slog.Warn("This is a warning message", "huba", "fjall")

	slog.Error("This is an error message", "huba", "fjall")

	slog.Debug("This is a debug message", "huba", "fjall", "password", "secret")

	spanProcessor := trace.NewSimpleSpanProcessor(&slogseq.LoggingSpanProcessor{Handler: handler})
	tp := trace.NewTracerProvider(trace.WithSpanProcessor(spanProcessor), trace.WithSampler(trace.AlwaysSample()))
	tracer := tp.Tracer("example-tracer")
	ctx := context.Background()
	spanCtx, span := tracer.Start(ctx, "operation")
	span.AddEvent("Starting operation", tr.WithAttributes(attribute.String("huba", "fjall")))
	slog.InfoContext(spanCtx, "This is a message with a span", "huba", "fjall")
	slog.WarnContext(spanCtx, "This is a warning message with a span", "huba", "fjall")
	time.Sleep(1 * time.Second)
	span.AddEvent("Doing some work")
	span.End()
	spanCtx, span = tracer.Start(context.WithValue(spanCtx, "start", time.Now()), "operation")
	slog.ErrorContext(spanCtx, "This is an error message with a span", "huba", "fjall")
	time.Sleep(100 * time.Millisecond)
	span.AddEvent("Doing some more work")
	slog.DebugContext(spanCtx, "This is a debug message with a span", "huba", "fjall", "password", "balle")
	span.AddEvent("This is an event with a span")
	time.Sleep(400 * time.Millisecond)
	span.AddEvent("Finishing operation")
	slog.InfoContext(spanCtx, "This is a message with a span", "huba", "fjall")
	span.End()
	spanCtx, span = tracer.Start(context.WithValue(spanCtx, "start", time.Now()), "operation")
	slog.WarnContext(spanCtx, "This is a warning message with a span", "huba", "fjall")
	span.AddEvent("Not finished yet")
	time.Sleep(1 * time.Second)
	slog.ErrorContext(spanCtx, "This is an error message with a span", "huba", "fjall")
	span.AddEvent("I think I'm done")
	span.RecordError(fmt.Errorf("this is an error"))
	span.End()
}
