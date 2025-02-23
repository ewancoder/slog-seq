package slogseq

import "time"

// Compact Log Event Format (CLEF) is a JSON-based log event format that Seq uses.
// https://clef-json.org
type CLEFEvent struct {
	Timestamp          time.Time              `json:"@t,omitzero,omitempty"`
	Message            string                 `json:"@m,omitempty"`
	Level              string                 `json:"@l"`
	Properties         map[string]interface{} `json:"-"`
	TraceID            string                 `json:"@tr,omitempty"`
	SpanID             string                 `json:"@sp,omitempty"`
	SpanStart          time.Time              `json:"@st,omitempty,omitzero"`
	ResourceAttributes map[string]interface{} `json:"@ra,omitempty,omitzero"`
	ParentSpanID       string                 `json:"@ps,omitempty"`
}
