package leds

import (
	"github.com/faiface/beep"
)

type Sink struct {
	stream beep.Streamer
}

func NewSink(stream beep.Streamer) *Sink {
	return &Sink{stream: stream}
}

func (s *Sink) Run() {
	var samples = make([][2]float64, 128)
	for {
		if _, ok := s.stream.Stream(samples); !ok {
			return
		}
	}
}
