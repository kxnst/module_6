package dsp

type EffectInfo struct {
	ParameterDefinitions []ParameterDefinition `json:"parameterDefinitions"`
	Slug                 string                `json:"slug"`
	Name                 string                `json:"name"`
	DspType              string                `json:"dspType"`
}
