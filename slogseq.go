package slogseq

import (
	"log/slog"
	"net/http"
	"time"
)

// SeqOption is an option to configure a Seq handler.
type SeqOption interface {
	apply(*SeqHandler) *SeqHandler
}

type seqOptionFunc func(*SeqHandler) *SeqHandler

func (f seqOptionFunc) apply(h *SeqHandler) *SeqHandler {
	return f(h)
}

// NewLogger creates a new Seq logger. seqURL is the URL of the Seq server.
// opts is a list of options to configure the Seq handler.
func NewLogger(seqURL string, opts ...SeqOption) (*slog.Logger, *SeqHandler) {
	handler := newSeqHandler(seqURL)
	for _, opt := range opts {
		handler = opt.apply(handler)
	}
	handler.start()
	return slog.New(handler), handler
}

// WithAPIKey sets the API key for the Seq server.
func WithAPIKey(apiKey string) SeqOption {
	return seqOptionFunc(func(h *SeqHandler) *SeqHandler {
		h.apiKey = apiKey
		return h
	})
}

// WithBatchSize sets the number of events to batch before sending to Seq.
func WithBatchSize(batchSize int) SeqOption {
	return seqOptionFunc(func(h *SeqHandler) *SeqHandler {
		h.batchSize = batchSize
		return h
	})
}

// WithFlushInterval sets the interval at which to flush the batch.
func WithFlushInterval(flushInterval time.Duration) SeqOption {
	return seqOptionFunc(func(h *SeqHandler) *SeqHandler {
		h.flushInterval = flushInterval
		return h
	})
}

// WithHandlerOptions sets the slog handler options.
func WithHandlerOptions(opts *slog.HandlerOptions) SeqOption {
	return seqOptionFunc(func(h *SeqHandler) *SeqHandler {
		h.options = *opts
		return h
	})
}

// WithInsecure disables TLS verification. Doesn't do anything if WithHTTPClient is also set.
func WithInsecure() SeqOption {
	return seqOptionFunc(func(h *SeqHandler) *SeqHandler {
		h.disableTLSVerify = true
		return h
	})
}

// WithHTTPClient sets the HTTP client to use. If not set, a default client is used.
func WithHTTPClient(client *http.Client) SeqOption {
	return seqOptionFunc(func(h *SeqHandler) *SeqHandler {
		h.client = client
		return h
	})
}

// WithGlobalAttrs sets the global attributes to include in all events.
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
