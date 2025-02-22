package slogseq

import (
	"log/slog"
	"time"
)

// NewSeqLogger creates a new Seq logger.
// seqURL is the URL of the Seq server.
// apiKey is the API key for the Seq server.
// batchSize is the number of events to batch before sending to Seq.
// flushInterval is the interval at which to flush the batch.
// opts is the handler options.
// Returns the logger and a function to close the logger.
//
// Example:
// 	seqLogger, finisher := slogseq.NewSeqLogger( ... )
// 	defer finisher()
// 	slog.SetDefault(seqLogger)
// 	slog.Info("Hello from slog-seq!")
func NewSeqLogger(seqURL, apiKey string, batchSize int, flushInterval time.Duration, opts *slog.HandlerOptions) (*slog.Logger, func() error) {
	handler := newSeqHandler(seqURL, apiKey, batchSize, flushInterval, opts)
	return slog.New(handler), handler.Close
}
