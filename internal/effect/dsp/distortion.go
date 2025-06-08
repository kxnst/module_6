package dsp

import (
	"guitar_processor/internal/context"
	"guitar_processor/internal/entity"
	"math"
)

type Distortion struct {
	ParameterStorage
	gainParam, levelParam *SliderParameter
	clipParam             *SwitchParameter
}

func (d *Distortion) GetSlug() string {
	return "distortion"
}

func tanhClip(x float32) float32 {
	return float32(math.Tanh(float64(x)))
}

func hardClip(x float32) float32 {
	if x > 1 {
		return 1
	} else if x < -1 {
		return -1
	}
	return x
}

func (d *Distortion) Process(samples []float32, _ *context.ProcessingContext) {
	for i := range samples {
		s := samples[i] * float32(math.Pow(10, float64(d.gainParam.GetValue())))

		if d.clipParam.GetValue() == SoftClip {
			s = hardClip(s)
		} else {
			s = tanhClip(s)
		}

		samples[i] = s * (d.levelParam.GetValue())
	}
}

const (
	SoftClip = 0
	HardClip = 1
)

func (d *Distortion) GetParameterDefinitions() []ParameterDefinition {
	return []ParameterDefinition{
		&SwitchParameterDefinition{
			Name:          "SoftClip",
			AllowedValues: map[string]float32{"soft": SoftClip, "hard": HardClip},
			Default:       SoftClip,
		},
		&SliderParameterDefinition{Name: "Gain", Min: 0, Max: 10, Default: 5},
		&SliderParameterDefinition{Name: "Level", Min: 0, Max: 1, Default: 0.7},
	}
}

func (d *Distortion) GetInfo(e *entity.Effect) EffectInfo {
	info := EffectInfo{
		Slug:                 d.GetSlug(),
		ParameterDefinitions: d.GetParameterDefinitions(),
	}
	if e != nil {
		info.DspType = e.DSPType
		info.Name = e.Name
	}
	return info
}

func (d *Distortion) Init(_ *context.ProcessingContext) {
	for _, param := range d.Parameters {
		switch def := param.GetDefinition().(type) {
		case *SliderParameterDefinition:
			switch def.Name {
			case "Gain":
				d.gainParam = param.(*SliderParameter)
			case "Level":
				d.levelParam = param.(*SliderParameter)
			}
		case *SwitchParameterDefinition:
			if def.Name == "SoftClip" {
				d.clipParam = param.(*SwitchParameter)
			}
		}
	}
}
