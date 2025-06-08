package request

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type DeviceInfoRequest struct {
	Device1Index int `json:"device1Index"`
	Device2Index int `json:"device2Index"`
	SampleRate   int `json:"sampleRate"`
	BufferSize   int `json:"bufferSize"`
}
