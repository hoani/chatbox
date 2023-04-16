package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/faiface/beep"
	"github.com/hoani/chatbox/3rdparty/faiface/beep/wav"
	"github.com/hoani/chatbox/hal"
	"github.com/hoani/chatbox/leds"
	"github.com/hoani/toot"
)

func (c *chatbox) doStateListening() state {
	start := time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	path := filepath.Join(c.wd, "test.wav")
	f, err := os.Create(path)
	if err != nil {
		c.errorMessage = "unable to record to file"
		c.hal.Debug(err.Error())
		return stateError
	}
	m, err := toot.NewDefaultMicrophone()
	if err != nil {
		c.errorMessage = "unable to open mic"
		c.hal.Debug(err.Error())
		return stateError
	}
	defer func() {
		go func() {
			time.Sleep(time.Second) // Delay a bit before stopping.
			m.Close()
		}()
	}()

	v := leds.NewVisualizer(
		leds.WithSource(&leds.Source{
			Streamer:   m,
			SampleRate: m.Format().SampleRate,
		}),
		leds.WithSink(func(s beep.Streamer) {
			go func() {
				err = wav.Encode(f, s, m.Format())
				if err != nil {
					c.hal.Debug(fmt.Sprintf("error encoding wav: %v", err))
					path = ""
				}
				f.Close()
				c.hal.Debug(path)
				c.recordingCh <- path
			}()
		}),
	)
	go v.Start(ctx)

	// c.hal.Debug(fmt.Sprintf("%#v\n", m.DeviceInfo()))

	if err := m.Start(ctx); err != nil {
		c.errorMessage = "unable to start mic"
		c.hal.Debug(err.Error())
		return stateError
	}

	c.hal.LCD().Write("  [Listening]  ", "release to stop", hal.LCDGreen)

	hsvs := []hal.HSV{}
	for i := 0; i < 24; i++ {
		hsvs = append(hsvs, hal.HSV{
			H: 0x60,
			S: 0xFF,
			V: 0x50,
		})
	}

	voicePowerEstimate := 0.0

	for {
		if time.Since(start) > 2*time.Minute {
			m.Close()
			c.errorMessage = "recording is too long"
			return stateError
		}
		if !c.hal.Button() {
			if time.Since(start) < buttonDebounce {
				continue // Allow for debounce.
			}
			if time.Since(start) < time.Second {
				m.Close()
				c.errorMessage = "recording is too short"
				return stateError
			}
			averagePowerEstimate := voicePowerEstimate / time.Since(start).Seconds()
			if averagePowerEstimate < 10.0 {
				m.Close()
				c.errorMessage = fmt.Sprintf("recording is too quiet %.2f", averagePowerEstimate)
				return stateError
			}
			break
		}

		time.Sleep(time.Millisecond * 20)

		channels := v.Channels()

		N := len(channels) / 4
		for i := 0; i < N; i++ {
			voicePowerEstimate += 0.02 * channels[i] / float64(N)
		}

		for i := range hsvs {
			j := i
			if j >= leds.NChannels {
				j = leds.NChannels - (1 + i - leds.NChannels)
			}
			v := channels[j]
			if v > float64(0xa0) {
				v = float64(0xa0)
			}
			hsvs[i].V = 0x40 + uint8(v)
		}

		c.hal.Leds().HSV(0, hsvs...)
		c.hal.Leds().Show()
	}

	return stateThinking
}
