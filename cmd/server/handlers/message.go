package handlers

import "guitar_processor/internal/effect/dsp"

type WebSocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type EffectAddedMessage struct {
	Type       string          `json:"type"`
	Slug       string          `json:"slug"`
	Parameters []dsp.Parameter `json:"parameters"`
}
