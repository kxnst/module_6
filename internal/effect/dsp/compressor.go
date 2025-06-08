package dsp

import (
	"guitar_processor/internal/context"
	"guitar_processor/internal/entity"
	"math"
)

type Compressor struct {
	ParameterStorage

	thresholdParam *SliderParameter
	ratioParam     *SliderParameter
	gainParam      *SliderParameter
}

func (c *Compressor) GetSlug() string {
	return "compressor"
}

func dbToLinear(db float32) float32 {
	return float32(math.Pow(10, float64(db)/20))
}

func linearToDb(linear float32) float32 {
	if linear <= 0.00001 {
		linear = 0.00001 // avoid log(0)
	}
	return float32(20 * math.Log10(float64(linear)))
}

func (c *Compressor) Process(samples []float32, _ *context.ProcessingContext) {
	thresholdDb := linearToDb(c.thresholdParam.GetValue())
	ratio := c.ratioParam.GetValue()
	makeupGain := dbToLinear(c.gainParam.GetValue())

	for i := range samples {
		input := samples[i]
		levelDb := linearToDb(abs(input))

		var output float32

		if levelDb > thresholdDb {
			excessDb := levelDb - thresholdDb
			gainReductionDb := excessDb - (excessDb / ratio)
			gain := dbToLinear(-gainReductionDb)
			output = input * gain
		} else {
			output = input
		}

		// Apply makeup gain and limit
		output *= makeupGain
		if output > 1 {
			output = 1
		} else if output < -1 {
			output = -1
		}

		samples[i] = output
	}
}

func (c *Compressor) GetParameterDefinitions() []ParameterDefinition {
	return []ParameterDefinition{
		&SliderParameterDefinition{
			Name:    "Threshold",
			Min:     -60,
			Max:     0,
			Default: -20,
		},
		&SliderParameterDefinition{
			Name:    "Ratio",
			Min:     1,
			Max:     20,
			Default: 4,
		},
		&SliderParameterDefinition{
			Name:    "MakeupGain",
			Min:     0,
			Max:     12,
			Default: 0,
		},
	}
}

func (c *Compressor) GetInfo(e *entity.Effect) EffectInfo {
	info := EffectInfo{
		Slug:                 c.GetSlug(),
		ParameterDefinitions: c.GetParameterDefinitions(),
	}
	if e != nil {
		info.DspType = e.DSPType
		info.Name = e.Name
	}
	return info
}

func (c *Compressor) Init(_ *context.ProcessingContext) {
	for _, param := range c.Parameters {
		switch def := param.GetDefinition().(type) {
		case *SliderParameterDefinition:
			switch def.Name {
			case "Threshold":
				c.thresholdParam = param.(*SliderParameter)
			case "Ratio":
				c.ratioParam = param.(*SliderParameter)
			case "MakeupGain":
				c.gainParam = param.(*SliderParameter)
			}
		}
	}
}

func abs(f float32) float32 {
	if f < 0 {
		return -f
	}
	return f
}
