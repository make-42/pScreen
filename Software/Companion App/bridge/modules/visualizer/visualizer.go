package visualizer

import (
	"image"
	"math"
	"math/cmplx"
	"pscreen/bridge/modules"
	"pscreen/bridge/renderer"
	"pscreen/config"
	"pscreen/utils"
	"time"

	"git.sr.ht/~sbinet/gg"
	"github.com/KarpelesLab/static-portaudio/portaudio"
	"github.com/mjibson/go-dsp/fft"
)

var firstFrame = true
var stream *portaudio.Stream
var signals = make(chan bool)
var easedMaxAbsSampleValue = 1.0
var freeSampleBufferL []float32
var freeSampleBufferR []float32
var Y []float64
var barCount = 0

func initMic() {
	var monoSampleBufferL = make([]float32, config.Config.VisualizerSampleBufferSize)
	var monoSampleBufferR = make([]float32, config.Config.VisualizerSampleBufferSize)
	freeSampleBufferL = make([]float32, config.Config.VisualizerCumulativeSampleBufferSize)
	freeSampleBufferR = make([]float32, config.Config.VisualizerCumulativeSampleBufferSize)
	portaudio.Initialize()
	var deviceOfInterest *portaudio.DeviceInfo
	/*devs, err := portaudio.Devices()
	utils.CheckError(err)

	/*fmt.Println("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
	for _, device := range devs {
		/*fmt.Println(device.Name, device.DefaultSampleRate, device.DefaultLowOutputLatency)
		if device.Name == "pipewire" /*"USB PnP Audio Device Mono" "Starship/Matisse HD Audio Controller Analog Stereo" {
			deviceOfInterest = device
		}
	}*/
	var err error
	if config.Config.VisualizerUseMicrophone {
		deviceOfInterest, err = portaudio.DefaultInputDevice()
	} else {
		deviceOfInterest, err = portaudio.DefaultOutputDevice()
	}
	utils.CheckError(err)
	streamParameters := portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device:   deviceOfInterest,
			Channels: 2,
			Latency:  time.Millisecond * time.Duration(config.Config.VisualizerInputDelayMillis),
		},
		SampleRate:      float64(config.Config.VisualizerSampleRate),
		FramesPerBuffer: config.Config.VisualizerSampleBufferSize,
	}
	stream, err = portaudio.OpenStream(streamParameters, monoSampleBufferL, monoSampleBufferR)
	utils.CheckError(err)
	utils.CheckError(stream.Start())
	for {
		utils.CheckError(stream.Read())
		freeSampleBufferL = freeSampleBufferL[config.Config.VisualizerSampleBufferSize:]
		freeSampleBufferR = freeSampleBufferR[config.Config.VisualizerSampleBufferSize:]
		freeSampleBufferL = append(freeSampleBufferL, monoSampleBufferL...)
		freeSampleBufferR = append(freeSampleBufferR, monoSampleBufferR...)
		select {
		case signals <- true:
			<-signals
		default:
			//
		}
	}
}

var VisualizerModule modules.Module = modules.Module{RenderFunction: func(im *image.RGBA) *image.RGBA {
	if firstFrame {
		firstFrame = false
		go initMic()
		barCount = config.Config.CanvasRenderDimensions.X / (config.Config.VisualizerFFTBarSpacing + config.Config.VisualizerFFTBarWidth)
		Y = make([]float64, barCount)
	}
	dc := gg.NewContextForRGBA(im)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetRGB(1, 1, 1)
	<-signals
	copiedFreeSampleBufferL := make([]float32, config.Config.VisualizerCumulativeSampleBufferSize)
	copiedFreeSampleBufferR := make([]float32, config.Config.VisualizerCumulativeSampleBufferSize)
	copy(copiedFreeSampleBufferL, freeSampleBufferL)
	copy(copiedFreeSampleBufferR, freeSampleBufferR)
	signals <- true
	maxAbsSampleValue := 0.0000000000001
	freeMergedSampleBuffer := make([]float64, config.Config.VisualizerCumulativeSampleBufferSize)
	for i := 0; i < config.Config.VisualizerCumulativeSampleBufferSize; i++ {
		freeMergedSampleBuffer[i] = float64(copiedFreeSampleBufferL[i]) + float64(copiedFreeSampleBufferR[i])
		if !config.Config.VisualizerShowFFT {
			if math.Abs(freeMergedSampleBuffer[i]) > maxAbsSampleValue {
				maxAbsSampleValue = math.Abs(freeMergedSampleBuffer[i])
			}
		}
	}

	maxPeriodLength := 0.0
	maxFreqPhaseFrac := 0.0
	X := fft.FFTReal(freeMergedSampleBuffer)
	if !config.Config.VisualizerShowFFT {
		maxFreq := 0
		maxFreqAmp := 0.0
		maxFreqPhase := 0.0

		for i := 0; i < config.Config.VisualizerCumulativeSampleBufferSize/4; i++ {
			freqAmpCmplx := X[i+config.Config.VisualizerCumulativeSampleBufferSize/2]
			freqAmp := math.Sqrt(real(cmplx.Conj(freqAmpCmplx) * freqAmpCmplx))
			if freqAmp > maxFreqAmp {
				maxFreq = i
				maxFreqAmp = freqAmp
				maxFreqPhase = math.Atan(imag(freqAmpCmplx) / real(freqAmpCmplx))
			}
		}
		maxFreqAdjusted := float64(maxFreq) * float64(config.Config.VisualizerSampleRate) / (float64(config.Config.VisualizerCumulativeSampleBufferSize) / 2)
		maxPeriodLength = float64(config.Config.VisualizerSampleRate) / maxFreqAdjusted
		maxFreqPhaseFrac = maxFreqPhase / (2 * math.Pi)
	} else {
		step := int(math.Round(float64(config.Config.VisualizerCumulativeSampleBufferSize) / float64(barCount) / 4 * float64(config.Config.VisualizerFFTCutoff)))
		for i := 0; i < barCount; i++ {
			dist := 0.0
			for j := 0; j < step; j++ {
				index := config.Config.VisualizerCumulativeSampleBufferSize/2 + i*step + j
				dist += math.Sqrt(math.Pow(real(X[index]), 2) + math.Pow(imag(X[index]), 2))
			}
			if maxAbsSampleValue < dist {
				maxAbsSampleValue = dist
			}
			Y[i] = dist*(1-config.Config.VisualizerFFTSmoothing) + Y[i]*config.Config.VisualizerFFTSmoothing
		}
	}

	easedMaxAbsSampleValue = easedMaxAbsSampleValue*config.Config.VisualizerScaleSmoothing + maxAbsSampleValue*(1-config.Config.VisualizerScaleSmoothing)
	if config.Config.VisualizerShowFFT {
		if easedMaxAbsSampleValue < config.Config.VisualizerMinScale {
			easedMaxAbsSampleValue = config.Config.VisualizerMinScale
		}
	}
	usedMaxAbsSampleValue := easedMaxAbsSampleValue / config.Config.VisualizerScale

	if config.Config.VisualizerShowFFT {
		for x := 0; x < barCount; x++ {
			dist := Y[x] / usedMaxAbsSampleValue
			bar_height := dist * float64(config.Config.CanvasRenderDimensions.Y)
			dc.DrawRectangle(float64(x*(config.Config.VisualizerFFTBarSpacing+config.Config.VisualizerFFTBarWidth)), float64(config.Config.CanvasRenderDimensions.Y)-bar_height, float64(config.Config.VisualizerFFTBarWidth), bar_height)
		}
	} else {
		offset := int(math.Round(maxPeriodLength * (1 - maxFreqPhaseFrac)))
		if math.Abs(float64(offset)) > float64(config.Config.VisualizerCumulativeSampleBufferSize) {
			offset = 0
		}
		samplesToShow := int(math.Max(math.Min(math.Round(maxPeriodLength), float64(config.Config.VisualizerCumulativeSampleBufferSize-offset)), float64(config.Config.CanvasRenderDimensions.X)))
		for x := 0; x < config.Config.CanvasRenderDimensions.X; x++ {
			dist := freeMergedSampleBuffer[x*samplesToShow/config.Config.CanvasRenderDimensions.X+offset]
			dc.LineTo(float64(x), (dist/usedMaxAbsSampleValue+1)*float64(config.Config.CanvasRenderDimensions.Y)/2)
		}
	}
	if config.Config.VisualizerShowFFT {
		dc.Fill()
	} else {
		dc.Stroke()
	}
	return renderer.RemoveAntiAliasing(im)
}}
