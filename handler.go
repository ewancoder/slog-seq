package slogseq

import (
	"context"
	"log/slog"
	"strings"
	"sync"
	"time"
)

type concurrencyState struct {
	eventsCh chan CLEFEvent
	doneCh   chan struct{}
	wg       sync.WaitGroup
}

type SeqHandler struct {
	// config
	seqURL        string
	apiKey        string
	batchSize     int
	flushInterval time.Duration

	// concurrency
	state *concurrencyState

	// Other fields for global attrs, grouping, etc.
	attrs  []slog.Attr
	groups []string
}

func NewSeqHandler(seqURL string, apiKey string, batchSize int, flushInterval time.Duration) *SeqHandler {
	h := &SeqHandler{
		seqURL:        seqURL,
		apiKey:        apiKey,
		batchSize:     batchSize,
		flushInterval: flushInterval,
		state: &concurrencyState{
			eventsCh: make(chan CLEFEvent, 1000), // some buffer size
			doneCh:   make(chan struct{}),
		},
	}

	// Start background flusher
	h.state.wg.Add(1)
	go h.runBackgroundFlusher()

	return h
}

func (h *SeqHandler) Handle(ctx context.Context, r slog.Record) error {
	// Convert slog.Level to text
	levelString := convertLevel(r.Level)

	// Collect attributes into a map
	props := make(map[string]interface{})
	h.addAttrs(props, h.attrs)
	r.Attrs(func(a slog.Attr) bool {
		h.addAttr(props, a)
		return true
	})

	// Create CLEF event
	event := CLEFEvent{
		Timestamp:  r.Time,
		Message:    r.Message,
		Level:      levelString,
		Properties: props,
	}

	// Send to channel (non-blocking or minimal blocking)
	select {
	case h.state.eventsCh <- event:
		// success
	default:
		// channel is full -> decide if we drop or block
		// for non-blocking, you might drop
		// or you can block (which might block the application)
		// or consider a better queue strategy
	}
	return nil
}

func (h *SeqHandler) Enabled(ctx context.Context, l slog.Level) bool {
	// This handler is always enabled
	return true
}

func (h *SeqHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h2 := *h
	h2.attrs = make([]slog.Attr, len(h.attrs))
	copy(h2.attrs, h.attrs)
	h2.attrs = append(h2.attrs, attrs...)

	return &h2
}

func (h *SeqHandler) WithGroup(name string) slog.Handler {
	h2 := *h
	h2.groups = make([]string, len(h.groups))
	copy(h2.groups, h.groups)
	h2.groups = append(h2.groups, name)

	return &h2
}

func (h *SeqHandler) Close() error {
	close(h.state.doneCh)
	h.state.wg.Wait()

	return nil
}

func (h *SeqHandler) addAttrs(dst map[string]any, attrs []slog.Attr) {
	for _, a := range attrs {
		h.addAttr(dst, a)
	}
}

func (h *SeqHandler) addAttr(dst map[string]any, a slog.Attr) {
	if a.Key == "" {
		return
	}
	var finalKey string
	if len(h.groups) > 0 {
		finalKey = strings.Join(h.groups, ".") + "." + a.Key
	} else {
		finalKey = a.Key
	}

	dst[finalKey] = a.Value.Any()
}

func convertLevel(l slog.Level) string {
	switch l {
	case slog.LevelDebug:
		return "Debug"
	case slog.LevelInfo:
		return "Information"
	case slog.LevelWarn:
		return "Warning"
	case slog.LevelError:
		return "Error"
	default:
		return "Unknown"
	}
}
