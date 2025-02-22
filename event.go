package slogseq

import "time"

type CLEFEvent struct {
	Timestamp  time.Time              `json:"@t"`
	Message    string                 `json:"@m,omitempty"`
	Level      string                 `json:"@l"`
	Properties map[string]interface{} `json:"-"`
}
