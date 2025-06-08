package dsp

import (
	"fmt"
	"github.com/google/uuid"
)

type ParameterStorage struct {
	Parameters []Parameter `json:"parameters"`
}

func (ps *ParameterStorage) InitParameters(params []ParameterDefinition) {
	var parameters []Parameter

	for _, definition := range params {
		if definition.GetType() == string(Switch) {
			parameter := &SwitchParameter{
				Definition: definition.(*SwitchParameterDefinition),
				Value:      definition.(*SwitchParameterDefinition).Default,
				UUID:       uuid.New().String(),
			}

			parameters = append(parameters, parameter)
		} else if definition.GetType() == string(Slider) {
			parameter := &SliderParameter{
				Definition: definition.(*SliderParameterDefinition),
				Value:      definition.(*SliderParameterDefinition).Default,
				UUID:       uuid.New().String(),
			}
			parameters = append(parameters, parameter)
		}
	}

	ps.Parameters = parameters
}

func (ps *ParameterStorage) GetParameterValue(uuid string) (float32, error) {
	for _, parameter := range ps.Parameters {
		if parameter.GetUUID() == uuid {
			return parameter.GetValue(), nil
		}
	}

	return 0, fmt.Errorf("parameter not found")
}

func (ps *ParameterStorage) SetParameterValue(uuid string, value float32) error {
	for _, parameter := range ps.Parameters {
		if parameter.GetUUID() == uuid {
			err := parameter.SetValue(value)

			return err
		}
	}

	return fmt.Errorf("parameter not found")
}

func (ps *ParameterStorage) GetParameters() []Parameter {
	return ps.Parameters

}
