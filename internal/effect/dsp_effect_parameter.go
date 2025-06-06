package effect

type DspEffectParameter interface {
	GetType() string
}

type SliderParameter struct {
	Name    string
	Min     float32
	Max     float32
	Default float32
}

type SwitchParameter struct {
	Name          string
	AllowedValues []float32
	ActiveValue   float32
}

func (s *SliderParameter) GetType() string {
	return "slider"
}

func (s *SwitchParameter) GetName() string {
	return "switch"
}

func (s *SwitchParameter) SetActiveValue(value float32) bool {
	for _, allowed := range s.AllowedValues {
		if allowed == value {
			s.ActiveValue = value
			return true
		}
	}

	return false
}
