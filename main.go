package main

import (
	"bufio"
	"fmt"
	"guitar_processor/internal/effect"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gordonklaus/portaudio"
)

func main() {
	must(portaudio.Initialize())
	defer portaudio.Terminate()

	devices, err := portaudio.Devices()
	must(err)

	// 🎙️ Крок 1: обираємо вхід
	fmt.Println("🔍 Список пристроїв з драйверами:")
	for i, dev := range devices {
		if dev.MaxInputChannels > 0 {
			fmt.Printf("[%d] %s [driver: %s]\n", i, dev.Name, dev.HostApi.Name)
		}
	}
	fmt.Print("🟢 Введіть індекс вхідного пристрою: ")
	inputDev := devices[24]

	// 🎧 Крок 2: обираємо вихід
	fmt.Println("\n🎧 Доступні ВИХІДНІ пристрої:")
	for i, dev := range devices {
		if dev.MaxOutputChannels > 0 {
			fmt.Printf("[%d] %s [driver: %s]\n", i, dev.Name, dev.HostApi.Name)
		}
	}
	fmt.Print("🔵 Введіть індекс вихідного пристрою: ")
	outputDev := devices[19]

	// ▶️ Крок 3: запускаємо потік
	fmt.Printf("\n▶️ Старт: %s → %s\n", inputDev.Name, outputDev.Name)
	stream := mustCreateStream(inputDev, outputDev)
	defer stream.Close()
	must(stream.Start())

	fmt.Println("🔊 Потік запущено. Натисни Ctrl+C для виходу.")
	select {}
}

func mustCreateStream(input, output *portaudio.DeviceInfo) *portaudio.Stream {
	const sampleRate = 44100
	const bufferSize = 256
	fmt.Println("Вхід має семплрейт:", input.DefaultSampleRate)
	fmt.Println("Вихід семплрейт:", output.DefaultSampleRate)
	//reverb := effect.NewReverb(sampleRate, 200.0, 0.3) // 80 мс затримка, 40% decay

	stream, err := portaudio.OpenStream(portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device:   input,
			Channels: input.MaxInputChannels,
			Latency:  input.DefaultHighInputLatency,
		},
		Output: portaudio.StreamDeviceParameters{
			Device:   output,
			Channels: output.MaxOutputChannels,
			Latency:  output.DefaultHighInputLatency,
		},
		SampleRate:      sampleRate,
		FramesPerBuffer: bufferSize,
	}, func(in, out []float32) {
		frames := len(in) / 2
		buffer := make([]float32, frames)

		for i := 0; i < frames; i++ {
			buffer[i] = in[i*2] // лише правий вхід
		}
		ds := &effect.BossDistortion{
			Gain:     8.0,
			Level:    0.9,
			HardClip: true, // або true для square-style
		}
		ds.Process(buffer)

		//reverb.Process(buffer)
		for i := 0; i < frames; i++ {
			s := buffer[i]
			out[i*4] = s   // лівий
			out[i*4+1] = s // правий
			out[i*4+2] = 0 // mute
			out[i*4+3] = 0
		}

	},
	)
	must(err)
	return stream
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func readInt() int {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		idx, err := strconv.Atoi(input)
		if err == nil {
			return idx
		}
		fmt.Print("❌ Введіть число: ")
	}
}

func clamp(x float32) float32 {
	if x > 1 {
		return 1
	}
	if x < -1 {
		return -1
	}
	return x
}
