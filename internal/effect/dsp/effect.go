package dsp

import (
	"guitar_processor/internal/context"
	"guitar_processor/internal/entity"
)

type Effect interface {
	Process(samples []float32, ctx *context.ProcessingContext)
	GetParameterDefinitions() []ParameterDefinition
	GetSlug() string
	GetInfo(e *entity.Effect) EffectInfo
	GetParameterValue(uuid string) (float32, error)
	SetParameterValue(uuid string, value float32) error
	Init(ctx *context.ProcessingContext)
	InitParameters(params []ParameterDefinition)
	GetParameters() []Parameter
}
