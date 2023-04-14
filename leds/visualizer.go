package leds

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/hoani/toot"
)

const NChannels = 12

type Source struct {
	Streamer   beep.Streamer
	SampleRate beep.SampleRate
}

type visualizer struct {
	channels     [NChannels]float64
	channelsLock sync.Mutex
	sink         func(s beep.Streamer)
	source       *Source
}

func NewVisualizer(options ...func(*visualizer)) *visualizer {
	v := &visualizer{}
	for _, o := range options {
		o(v)
	}
	if v.sink == nil {
		v.sink = func(s beep.Streamer) {
			d := NewSink(s)
			go d.Run()
		}
	}
	return v
}

func WithSink(sink func(s beep.Streamer)) func(*visualizer) {
	return func(v *visualizer) {
		v.sink = sink
	}
}

func WithSource(source *Source) func(*visualizer) {
	return func(v *visualizer) {
		v.source = source
	}
}

func (v *visualizer) Start(ctx context.Context) error {
	if v.source == nil {
		m, err := toot.NewDefaultMicrophone()
		if err != nil {
			return err
		}
		defer m.Close()
		v.source = &Source{
			Streamer:   m,
			SampleRate: m.Format().SampleRate,
		}
		go m.Start(ctx)
	}

	a := toot.NewAnalyzer(v.source.Streamer, int(v.source.SampleRate), int(v.source.SampleRate/4))
	tv := toot.NewVisualizer(100.0, 4000.0, NChannels)

	v.sink(a)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(50 * time.Millisecond):
			s := a.GetPowerSpectrum()
			if s == nil {
				continue
			}

			result := tv.Bin(s)
			v.channelsLock.Lock()
			for i, r := range result {
				v.channels[i] = math.Log10(1+r*1000) * 5000 // Do some log scaling to make the power spectra show up nicer.
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
