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

	handler := slogseq.NewSeqHandler(
		*seqURL,       // seqURL
		*apiKey,       // apiKey
		50,            // batchSize
		2*time.Second, // flushInterval
		opts,          // opts
	)
	defer handler.Close()

	logger := slog.New(handler).With("app", "slog-seq").With("env", "dev").With("version", "1.0.0")

	logger.Info("Hello from slog-seq test command!",
		"env", "dev",
		"version", "1.0.0")

	logger.Warn("This is a warning message", "huba", "fjall")

	logger.Error("This is an error message", "huba", "fjall")

	logger.Debug("This is a debug message", "huba", "fjall", "password", "secret")
}
