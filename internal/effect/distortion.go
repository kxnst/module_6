package effect

type Distortion struct {
	Gain     float32
	Level    float32
	HardClip bool
}

func softClip(x float32) float32 {
	if x < -1 {
		return -1
	}
	if x > 1 {
		return 1
	}
	return x - (x * x * x / 3)
}

func hardClip(x float32) float32 {
	if x > 1 {
		return 1
	} else if x < -1 {
		return -1
	}
	return x
}

func (d *Distortion) Process(samples []float32) {
	for i := range samples {
		s := samples[i] * d.Gain

		if d.HardClip {
			s = hardClip(s)
		} else {
			s = softClip(s)
		}

		samples[i] = s * d.Level
	}
}
