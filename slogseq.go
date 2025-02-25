package slogseq

import (
	"log/slog"
	"net/http"
	"time"
)

type SeqOption interface {
	apply(*SeqHandler) *SeqHandler
}

type seqOptionFunc func(*SeqHandler) *SeqHandler

func (f seqOptionFunc) apply(h *SeqHandler) *SeqHandler {
	return f(h)
}

func NewLogger(seqURL string, opts ...SeqOption) (*slog.Logger, *SeqHandler) {
	handler := newSeqHandler(seqURL)
	for _, opt := range opts {
		handler = opt.apply(handler)
	}
	handler.start()
	return slog.New(handler), handler
}

func WithAPIKey(apiKey string) SeqOption {
	return seqOptionFunc(func(h *SeqHandler) *SeqHandler {
		h.apiKey = apiKey
		return h
	})
}

func WithBatchSize(batchSize int) SeqOption {
	return seqOptionFunc(func(h *SeqHandler) *SeqHandler {
		h.batchSize = batchSize
		return h
	})
}

func WithFlushInterval(flushInterval time.Duration) SeqOption {
	return seqOptionFunc(func(h *SeqHandler) *SeqHandler {
		h.flushInterval = flushInterval
		return h
	})
}

func WithHandlerOptions(opts *slog.HandlerOptions) SeqOption {
	return seqOptionFunc(func(h *SeqHandler) *SeqHandler {
		h.options = *opts
		return h
	})
}

func WithInsecure() SeqOption {
	return seqOptionFunc(func(h *SeqHandler) *SeqHandler {
		h.disableTLSVerify = true
		return h
	})
}

func WithHTTPClient(client *http.Client) SeqOption {
	return seqOptionFunc(func(h *SeqHandler) *SeqHandler {
		h.client = client
		return h
	})
}

func WithGlobalAttrs(attrs ...slog.Attr) SeqOption {
	return seqOptionFunc(func(h *SeqHandler) *SeqHandler {
		h.attrs = attrs
		return h
	})
}

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

// Deprecated: Use NewLogger instead.
func NewSeqLogger(seqURL, apiKey string, batchSize int, flushInterval time.Duration, opts *slog.HandlerOptions) (*slog.Logger, *SeqHandler) {
	return NewLogger(seqURL,
		WithAPIKey(apiKey),
		WithBatchSize(batchSize),
		WithFlushInterval(flushInterval),
		WithHandlerOptions(opts),
	)
}
