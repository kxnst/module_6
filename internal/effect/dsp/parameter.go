package dsp

import "fmt"

type Parameter interface {
	GetUUID() string
	SetValue(value float32) error
	GetValue() float32
	GetDefinition() ParameterDefinition
}

type SliderParameter struct {
	UUID       string                     `json:"uuid"`
	Definition *SliderParameterDefinition `json:"definition"`
	Value      float32                    `json:"value"`
}

type SwitchParameter struct {
	UUID       string                     `json:"uuid"`
	Definition *SwitchParameterDefinition `json:"definition"`
	Value      float32                    `json:"value"`
}

func (s *SliderParameter) SetValue(value float32) error {
	if value < s.Definition.Min || value > s.Definition.Max {
		return fmt.Errorf("value %s not supported", fmt.Sprint(value))
	}
	s.Value = value
	return nil
}

func (s *SwitchParameter) SetValue(value float32) error {
	for _, v := range s.Definition.AllowedValues {
		if v == value {
			s.Value = value
			return nil
		}
	}

	return fmt.Errorf("option %s not supported", fmt.Sprint(value))
}

func (s *SwitchParameter) GetValue() float32 {
	return s.Value
}

func (s *SliderParameter) GetValue() float32 {
	return s.Value
}

func (s *SwitchParameter) GetUUID() string {
	return s.UUID
}

func (s *SliderParameter) GetUUID() string {
	return s.UUID
}

func (s *SwitchParameter) GetDefinition() ParameterDefinition {
	return s.Definition
}

func (s *SliderParameter) GetDefinition() ParameterDefinition {
	return s.Definition
}
