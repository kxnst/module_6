package dsp

const (
	DeviceInput  = 0
	DeviceOutput = 1
)

type Device struct {
	Name        string  `json:"name"`
	HostApiName string  `json:"hostApiName"`
	SampleRate  float64 `json:"sampleRate"`
	Type        int     `json:"type"`
	Channels    int     `json:"channels"`
	Index       int     `json:"index"`
}
