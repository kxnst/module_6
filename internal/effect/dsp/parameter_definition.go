package dsp

import "encoding/json"

type ParameterDefinition interface {
	GetType() string
}

type SliderParameterDefinition struct {
	Name    string  `json:"name"`
	Min     float32 `json:"min"`
	Max     float32 `json:"max"`
	Default float32 `json:"default"`
}

type SwitchParameterDefinition struct {
	Name          string             `json:"name"`
	AllowedValues map[string]float32 `json:"allowedValues"`
	Default       float32            `json:"default"`
}

type ParameterType string

const (
	Slider ParameterType = "slider"
	Switch ParameterType = "switch"
)

func (s *SliderParameterDefinition) GetType() string {
	return string(Slider)
}

func (s *SwitchParameterDefinition) GetType() string {
	return string(Switch)
}

// MarshalJSON customizes the JSON marshaling by adding a "type" field.
func (s *SwitchParameterDefinition) MarshalJSON() ([]byte, error) {
	type Alias SwitchParameterDefinition
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  s.GetType(),
		Alias: (*Alias)(s),
	})
}

func (s *SliderParameterDefinition) MarshalJSON() ([]byte, error) {
	type Alias SliderParameterDefinition
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  s.GetType(),
		Alias: (*Alias)(s),
	})
}
