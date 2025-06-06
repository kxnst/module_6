package effect

type Reverb struct {
	delayBuffer  []float32
	writeIndex   int
	delaySamples int
	decay        float32
}

func NewReverb(sampleRate float64, delayMs float64, decay float32) *Reverb {
	delaySamples := int(sampleRate * delayMs / 1000)
	return &Reverb{
		delayBuffer:  make([]float32, delaySamples),
		writeIndex:   0,
		delaySamples: delaySamples,
		decay:        decay,
	}
}

func (r *Reverb) Process(samples []float32) {
	for i := 0; i < len(samples); i++ {
		delayed := r.delayBuffer[r.writeIndex]
		dry := samples[i]
		wet := dry + delayed*r.decay

		r.delayBuffer[r.writeIndex] = wet
		samples[i] = wet // або: samples[i] = dry*(1.0 - mix) + wet*mix

		r.writeIndex = (r.writeIndex + 1) % r.delaySamples
	}
}
