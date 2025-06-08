package service

import (
	"fmt"
	"guitar_processor/internal/context"
	"guitar_processor/internal/effect/dsp"
	"guitar_processor/internal/entity"
	"guitar_processor/internal/repository"
	"slices"
)

type EffectsRegistry struct {
	DspEffects []dsp.Effect
	repo       *repository.EffectRepository
}

func NewEffectsRegistry(repo *repository.EffectRepository) *EffectsRegistry {
	effects := make([]dsp.Effect, 0)
	effects = append(effects, &dsp.Distortion{})
	effects = append(effects, &dsp.NoiseGate{})
	effects = append(effects, &dsp.Reverb{})
	effects = append(effects, &dsp.Compressor{})

	return &EffectsRegistry{repo: repo, DspEffects: effects}
}

func (r *EffectsRegistry) GetEffects() []dsp.Effect {
	return r.DspEffects
}

func (r *EffectsRegistry) GetEffectsInfo() []*dsp.EffectInfo {
	dbEffects, err := r.repo.GetEffects()
	if err != nil {
		dbEffects = []*entity.Effect{}
	}

	effectsInfo := make([]*dsp.EffectInfo, 0)
	for _, effect := range r.DspEffects {
		id := slices.IndexFunc(dbEffects, func(e *entity.Effect) bool { return e.Slug == effect.GetSlug() })
		var itemp *entity.Effect
		if id == -1 {
			itemp = &entity.Effect{}
		} else {
			itemp = dbEffects[id]
		}

		info := effect.GetInfo(itemp)

		effectsInfo = append(effectsInfo, &info)
	}

	return effectsInfo
}

func (r *EffectsRegistry) CreateBySlug(slug string, ctx *context.ProcessingContext) (dsp.Effect, error) {
	for _, effect := range r.DspEffects {
		if effect.GetSlug() == slug {
			effect.InitParameters(effect.GetParameterDefinitions())
			effect.Init(ctx)

			return effect, nil
		}
	}

	return nil, fmt.Errorf("effect %s not found", slug)
}
