package effect

type Effect interface {
	Process(samples []float32)
}
