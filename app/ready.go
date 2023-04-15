package app

import (
	"context"
	"time"

	"github.com/hoani/chatbox/hal"
	"github.com/hoani/chatbox/leds"
)

func (c *chatbox) doStateReady() state {
	c.hal.LCD().Write("Press to start", "", hal.LCDBlue)

	if c.hal.Button() {
		time.Sleep(time.Millisecond * 10)
	}
	time.Sleep(time.Millisecond * 200) // We delay a little bit extra to allow for button debounce.

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hsvs := []hal.HSV{}
	for i := 0; i < 24; i++ {
		hsvs = append(hsvs, hal.HSV{
			H: uint8(i) * 10,
			S: 0xFF,
			V: 0x50,
		})
	}

	v := leds.NewVisualizer()
	go func() {
		if err := v.Start(ctx); err != nil {
			panic(err)
		}
	}()
	defer v.Wait()

	for {
		if c.hal.Button() {
			cancel()
			break
		}

		time.Sleep(20 * time.Millisecond)

		channels := v.Channels()

		for i := range hsvs {
			hsvs[i].H += 1
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
	return stateListening
}
