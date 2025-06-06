package effect

type DspEffect interface {
	Process(samples []float32)
	GetEffectParameters() []DspEffect
	GetName() string
}
