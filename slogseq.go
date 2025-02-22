package slogseq

import (
	"log/slog"
	"time"
)

func NewSeqLogger(seqURL, apiKey string, batchSize int, flushInterval time.Duration) *slog.Logger {
    handler := NewSeqHandler(seqURL, apiKey, batchSize, flushInterval)
    return slog.New(handler)
}
