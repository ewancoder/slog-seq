package main

import (
	"flag"
	"log/slog"
	"time"

	slogseq "github.com/sokkalf/slog-seq" // import your library
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
}
