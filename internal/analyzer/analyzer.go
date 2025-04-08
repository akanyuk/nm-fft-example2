package analyzer

import (
	"gonum.org/v1/gonum/dsp/fourier"
	"math"
	"math/cmplx"
)

type Analyzer struct {
	buffer *[]float32
	data   []float64
	fft    *fourier.FFT
	Bands  []uint16
}

func NewAnalyzer(buffer []float32) *Analyzer {
	size := len(buffer)
	return &Analyzer{
		fft:    fourier.NewFFT(size),
		buffer: &buffer,
		data:   make([]float64, size),
		Bands:  make([]uint16, 32),
	}
}

func (a *Analyzer) Process() []uint16 {
	for i, v := range *a.buffer {
		a.data[i] = float64(v)
	}

	coeff := a.fft.Coefficients(nil, a.data)
	magnitude := make([]float64, len(coeff)/2)
	for i := range magnitude {
		magnitude[i] = cmplx.Abs(coeff[i])
	}
	bandSize := len(magnitude) / len(a.Bands)

	for i := range a.Bands {
		start := i * bandSize
		end := start + bandSize
		if i == 31 {
			end = len(magnitude)
		}

		maxVal := 0.0
		for j := start; j < end; j++ {
			if magnitude[j] > maxVal {
				maxVal = magnitude[j]
			}
		}

		if maxVal > 0 {
			dbVal := 20 * math.Log10(maxVal)
			normalized := min(max((dbVal+10)*8, 0), 255)
			a.Bands[i] = uint16(normalized)
		}
	}

	return a.Bands
}
