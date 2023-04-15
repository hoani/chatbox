package app

import (
	"context"
	"time"

	"github.com/hoani/chatbox/hal"
	"github.com/hoani/chatbox/lcd"
	"github.com/hoani/chatbox/leds"
)

func (c *chatbox) doStateReady() state {
	c.readyLCD()

	// Wait for button release.
	for c.hal.Button() {
		time.Sleep(buttonDebounce)
	}
	time.Sleep(buttonDebounce) // We delay a little bit extra to allow for button debounce.

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
	go v.Start(ctx)
	defer v.Wait()

	for {
		if c.hal.Button() {
			cancel()
			break
		}

		time.Sleep(20 * time.Millisecond)

		c.readyLCD()
		c.readyLEDs(hsvs, v)
		c.restartOldChat()
	}
	return stateListening
}

func (c *chatbox) readyLCD() {
	timestamp := time.Now().Format(time.Kitchen)
	if len(timestamp) == len(time.Kitchen) {
		timestamp = " " + timestamp
	}
	c.hal.LCD().Write(lcd.Pad(timestamp), lcd.Pad("Press to start"), hal.LCDBlue)
}

func (c *chatbox) readyLEDs(hsvs []hal.HSV, v leds.Visualizer) {
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

// Restarts chat at midnight if no one is actively using the chatbox.
func (c *chatbox) restartOldChat() {
	now := time.Now()
	if now.Hour() == 0 && now.Minute() == 0 && now.Second() == 0 {
		if time.Since(c.lastChat) > time.Hour {
			c.chat = c.newChatRequest()
		}
	}
}
