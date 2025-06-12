package dsp

import (
	"github.com/gofiber/fiber/v2/log"
	"guitar_processor/internal/context"
	"guitar_processor/internal/entity"
)

type Reverb struct {
	ParameterStorage
	delayBuffer  []float32
	writeIndex   int
	delaySamples int
	lastRate     int
	initialized  bool
}

func (r *Reverb) GetSlug() string {
	return "reverb"
}

func (r *Reverb) GetParameterDefinitions() []ParameterDefinition {
	return []ParameterDefinition{
		&SliderParameterDefinition{Name: "Delay", Min: 10, Max: 1000, Default: 100},
		&SliderParameterDefinition{Name: "Decay", Min: 0.0, Max: 1.0, Default: 0.5},
	}
}

func (r *Reverb) initialize(ctx *context.ProcessingContext) error {
	delayUUID := r.Parameters[0].GetUUID()
	decayUUID := r.Parameters[1].GetUUID()

	delayMs, err := r.GetParameterValue(delayUUID)
	if err != nil {
		return err
	}
	decay, err := r.GetParameterValue(decayUUID)
	if err != nil {
		return err
	}

	r.delaySamples = int(float32(ctx.SampleRate) * float32(delayMs) / 1000.0)
	r.delayBuffer = make([]float32, r.delaySamples)
	r.writeIndex = 0
	r.lastRate = ctx.SampleRate
	r.initialized = true

	r.ParameterStorage.Parameters[1].SetValue(decay)

	return nil
}

func (r *Reverb) Process(samples []float32, _ *context.ProcessingContext) {
	decay, _ := r.GetParameterValue(r.Parameters[1].GetUUID())

	for i := 0; i < len(samples); i++ {
		delayed := r.delayBuffer[r.writeIndex]
		dry := samples[i]
		wet := dry + delayed*decay

		r.delayBuffer[r.writeIndex] = wet
		samples[i] = wet

		r.writeIndex = (r.writeIndex + 1) % r.delaySamples
	}
}

func (r *Reverb) GetInfo(e *entity.Effect) EffectInfo {
	info := EffectInfo{
		Slug:                 r.GetSlug(),
		ParameterDefinitions: r.GetParameterDefinitions(),
	}
	if e == nil {
		return info
	}

	info.DspType = e.DSPType
	info.Name = e.Name

	return info
}

func (r *Reverb) Init(ctx *context.ProcessingContext) {
	err := r.initialize(ctx)
	if err != nil {
		log.Warn("Reverb init failed:", err)
		return
	}
}
