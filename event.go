package slogseq

import "time"

// Compact Log Event Format (CLEF) is a JSON-based log event format that Seq uses.
// https://clef-json.org
type CLEFEvent struct {
	Timestamp  time.Time              `json:"@t"`
	Message    string                 `json:"@m,omitempty"`
	Level      string                 `json:"@l"`
	Properties map[string]interface{} `json:"-"`
}
