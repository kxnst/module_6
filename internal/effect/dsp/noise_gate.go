package dsp

import (
	"guitar_processor/internal/context"
	"guitar_processor/internal/entity"
	"math"
)

type NoiseGate struct {
	ParameterStorage
	thresholdParam *SliderParameter
}

func (ng *NoiseGate) GetSlug() string {
	return "noise_gate"
}

func (ng *NoiseGate) Process(samples []float32, _ *context.ProcessingContext) {
	threshold := ng.thresholdParam.GetValue()

	for i := range samples {
		if math.Abs(float64(samples[i])) < float64(threshold) {
			samples[i] = 0
		}
	}
}

func (ng *NoiseGate) GetParameterDefinitions() []ParameterDefinition {
	return []ParameterDefinition{
		&SliderParameterDefinition{
			Name:    "Threshold",
			Min:     0.001,
			Max:     0.5,
			Default: 0.02,
		},
	}
}

func (ng *NoiseGate) GetInfo(e *entity.Effect) EffectInfo {
	info := EffectInfo{
		Slug:                 ng.GetSlug(),
		ParameterDefinitions: ng.GetParameterDefinitions(),
	}
	if e != nil {
		info.DspType = e.DSPType
		info.Name = e.Name
	}
	return info
}

func (ng *NoiseGate) Init(_ *context.ProcessingContext) {
	for _, param := range ng.Parameters {
		if def, ok := param.GetDefinition().(*SliderParameterDefinition); ok && def.Name == "Threshold" {
			ng.thresholdParam = param.(*SliderParameter)
		}
	}
}
