package service

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gordonklaus/portaudio"
	"guitar_processor/cmd/server/request"
	"guitar_processor/internal/context"
	"guitar_processor/internal/effect/dsp"
	"sync"
)

type AudioService interface {
	GetDevices() ([]*dsp.Device, error)
	Start(request *request.DeviceInfoRequest) (Stream, error)
	GetEffectsInfo() []*dsp.EffectInfo
	AddEffectToChain(slug string) (dsp.Effect, error)
	SetEffectParameter(uuid string, value float32) error
	Terminate()
}

type NativeAudioService struct {
	mu sync.Mutex
	er *EffectsRegistry
	ec *EffectsChain
}

func NewAudioService(er *EffectsRegistry) AudioService {
	return &NativeAudioService{er: er, ec: &EffectsChain{Chain: make([]dsp.Effect, 0)}}
}

func (s *NativeAudioService) GetDevices() ([]*dsp.Device, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := portaudio.Initialize()
	defer portaudio.Terminate()

	if err != nil {
		return nil, err
	}

	osDevices, err := portaudio.Devices()
	if err != nil {
		return nil, err
	}

	var devices []*dsp.Device

	for index, osDevice := range osDevices {
		device := &dsp.Device{
			Name:        osDevice.Name,
			HostApiName: osDevice.HostApi.Name,
			SampleRate:  osDevice.DefaultSampleRate,
			Index:       index,
		}
		if osDevice.MaxInputChannels > 0 {
			device.Type = dsp.DeviceInput
			device.Channels = osDevice.MaxInputChannels
		} else {
			device.Type = dsp.DeviceOutput
			device.Channels = osDevice.MaxOutputChannels
		}

		devices = append(devices, device)
	}

	return devices, nil
}

func (s *NativeAudioService) Start(request *request.DeviceInfoRequest) (Stream, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := portaudio.Initialize()

	if err != nil {
		return nil, err
	}

	osDevices, err := portaudio.Devices()
	if err != nil {
		return nil, err
	}

	s.ec = &EffectsChain{Chain: make([]dsp.Effect, 0)}
	ctx := &context.ProcessingContext{SampleRate: request.SampleRate, BufferSize: request.BufferSize}
	s.ec.Ctx = ctx

	input := osDevices[request.Device1Index]
	output := osDevices[request.Device2Index]

	if (input == nil) || (output == nil) {
		return nil, fmt.Errorf("device not found")
	}

	log.Info("Starting audio processing")
	stream, err := portaudio.OpenStream(portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device:   input,
			Channels: input.MaxInputChannels,
			Latency:  input.DefaultLowInputLatency,
		},
		Output: portaudio.StreamDeviceParameters{
			Device:   output,
			Channels: output.MaxOutputChannels,
			Latency:  output.DefaultLowOutputLatency,
		},
		SampleRate:      float64(ctx.SampleRate),
		FramesPerBuffer: ctx.BufferSize,
	}, func(in, out []float32) {
		inputChannels := input.MaxInputChannels
		outputChannels := output.MaxOutputChannels
		frames := len(in) / inputChannels

		buffer := make([]float32, frames)

		for i := 0; i < frames; i++ {
			var mono float32
			for ch := 0; ch < inputChannels; ch++ {
				mono += in[i*inputChannels+ch]
			}
			buffer[i] = mono / float32(inputChannels)
		}

		for _, effect := range s.ec.Chain {
			effect.Process(buffer, ctx)
		}

		for i := 0; i < frames; i++ {
			sample := buffer[i]
			for ch := 0; ch < outputChannels; ch++ {
				out[i*outputChannels+ch] = sample
			}
		}
	},
	)

	return stream, err
}

func (s *NativeAudioService) GetEffectsInfo() []*dsp.EffectInfo {
	return s.er.GetEffectsInfo()
}

func (s *NativeAudioService) AddEffectToChain(slug string) (dsp.Effect, error) {
	effect, err := s.er.CreateBySlug(slug, s.ec.Ctx)
	if err != nil {
		return nil, err
	}
	s.ec.Chain = append(s.ec.Chain, effect)

	return effect, nil
}

func (s *NativeAudioService) SetEffectParameter(uuid string, value float32) error {
	return s.ec.SetParameter(uuid, value)
}

func (s *NativeAudioService) Terminate() {
	portaudio.Terminate()
}
