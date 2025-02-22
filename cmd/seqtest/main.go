package main

import (
	"context"
	"flag"
	"log/slog"
	"path"
	"time"

	slogseq "github.com/sokkalf/slog-seq" // import your library
	"go.opentelemetry.io/otel/sdk/trace"
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

	seqLogger, finisher := slogseq.NewSeqLogger(
		*seqURL,       // seqURL
		*apiKey,       // apiKey
		50,            // batchSize
		2*time.Second, // flushInterval
		opts,          // opts
	)
	defer finisher()

	slog.SetDefault(seqLogger.With("app", "slog-seq").With("env", "dev").With("version", "1.0.0"))

	slog.Info("Hello from slog-seq test command!",
		"env", "dev",
		"version", "1.0.0")

	slog.Warn("This is a warning message", "huba", "fjall")

	slog.Error("This is an error message", "huba", "fjall")

	slog.Debug("This is a debug message", "huba", "fjall", "password", "secret")

	tp := trace.NewTracerProvider(trace.WithSampler(trace.AlwaysSample()))
	tracer := tp.Tracer("example-tracer")
	ctx := context.Background()
	spanCtx, span := tracer.Start(context.WithValue(ctx, "start", time.Now()), "operation")

	slog.InfoContext(spanCtx, "This is a message with a span", "huba", "fjall")
	slog.WarnContext(spanCtx, "This is a warning message with a span", "huba", "fjall")
	time.Sleep(1 * time.Second)
	span.End()
	spanCtx, span = tracer.Start(context.WithValue(spanCtx, "start", time.Now()), "operation")
	slog.ErrorContext(spanCtx, "This is an error message with a span", "huba", "fjall")
	time.Sleep(100 * time.Millisecond)
	slog.DebugContext(spanCtx, "This is a debug message with a span", "huba", "fjall", "password", "balle")
	time.Sleep(400 * time.Millisecond)
	slog.InfoContext(spanCtx, "This is a message with a span", "huba", "fjall")
	span.End()
	spanCtx, span = tracer.Start(context.WithValue(spanCtx, "start", time.Now()), "operation")
	slog.WarnContext(spanCtx, "This is a warning message with a span", "huba", "fjall")
	time.Sleep(1 * time.Second)
	slog.ErrorContext(spanCtx, "This is an error message with a span", "huba", "fjall")
	span.End()
}
