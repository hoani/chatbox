package leds

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/hoani/toot"
)

const NChannels = 12

type visualizer struct {
	channels     [NChannels]float64
	channelsLock sync.Mutex
}

func NewVisualizer() *visualizer {
	return &visualizer{}
}

func (v *visualizer) Start(ctx context.Context) error {
	m, err := toot.NewDefaultMicrophone()
	if err != nil {
		return err
	}
	defer m.Close()

	a := toot.NewAnalyzer(m, int(m.Format().SampleRate), int(m.Format().SampleRate/4))
	tv := toot.NewVisualizer(100.0, 4000.0, NChannels)
	s := NewSink(a)

	go m.Start(ctx)
	go s.Run()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(time.Millisecond):
			s := a.GetPowerSpectrum()
			if s == nil {
				continue
			}

			result := tv.Bin(s)
			v.channelsLock.Lock()
			for i, r := range result {
				v.channels[i] = math.Log10(1+r*1000) * 1000 // Do some log scaling to make the power spectra show up nicer.
			}
			v.channelsLock.Unlock()
		}
	}
}

func (v *visualizer) Channels() [NChannels]float64 {
	v.channelsLock.Lock()
	defer v.channelsLock.Unlock()
	return v.channels
}
