package service

import (
	"fmt"
	"guitar_processor/internal/context"
	"guitar_processor/internal/effect/dsp"
)

type EffectsChain struct {
	Ctx   *context.ProcessingContext
	Chain []dsp.Effect
}

func (ec *EffectsChain) Add(effect dsp.Effect) {
	ec.Chain = append(ec.Chain, effect)
}

func (ec *EffectsChain) Flush() {
	ec.Chain = make([]dsp.Effect, 0)
}

func (ec *EffectsChain) SetParameter(uuid string, value float32) error {
	for _, eff := range ec.Chain {
		if err := eff.SetParameterValue(uuid, value); err == nil {
			return nil
		}
	}

	return fmt.Errorf("parameter %s not found", uuid)
}
