package unit

import (
	"guitar_processor/internal/context"
	"guitar_processor/internal/effect/dsp"
	"guitar_processor/internal/effect/service"
	"testing"
)

func TestEffectsChain_AddSetFlush(t *testing.T) {
	ctx := &context.ProcessingContext{SampleRate: 44100, BufferSize: 32}
	chain := &service.EffectsChain{Ctx: ctx}

	dist := &dsp.Distortion{}
	dist.InitParameters(dist.GetParameterDefinitions())
	dist.Init(ctx)

	var targetUUID = dist.GetParameters()[0].GetUUID()
	var targetVal float32
	for _, p := range dist.GetParameters() {
		if mapped, ok := p.GetDefinition().(*dsp.SliderParameterDefinition); ok {
			targetUUID = p.GetUUID()
			targetVal = (mapped.Min + mapped.Max) / 2
			break
		}
	}

	if targetUUID == "" {
		t.Fatal("Gain parameter not found")
	}

	chain.Add(dist)

	err := chain.SetParameter(targetUUID, targetVal)
	if err != nil {
		t.Errorf("failed to set parameter: %v", err)
	}

	for _, p := range dist.GetParameters() {
		if p.GetUUID() == targetUUID {
			if p.GetValue() != targetVal {
				t.Errorf("parameter value not set correctly, expected %v got %v", targetVal, p.GetValue())
			}
		}
	}

	chain.Flush()
	if len(chain.Chain) != 0 {
		t.Error("Flush() did not clear the effects chain")
	}
}
