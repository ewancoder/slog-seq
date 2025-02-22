package slogseq

import (
	"log/slog"
	"time"
)

func NewSeqLogger(seqURL, apiKey string, batchSize int, flushInterval time.Duration, opts *slog.HandlerOptions) *slog.Logger {
	handler := NewSeqHandler(seqURL, apiKey, batchSize, flushInterval, opts)
	return slog.New(handler)
}
