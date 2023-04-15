package app

import (
	"context"
	"math"
	"time"

	"github.com/hoani/chatbox/hal"
)

func (c *chatbox) doStateError() state {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		hsvs := []hal.HSV{}
		for i := 0; i < 24; i++ {
			hsvs = append(hsvs, hal.HSV{
				H: 0x00,
				S: 0xFF,
				V: 20 + 20*uint8(i),
			})
		}
		start := time.Now()
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Millisecond * 50):
				v := 0x80 + uint8(0x60*math.Sin(math.Pi*float64(time.Since(start).Seconds()/10)))
				for i := 0; i < 24; i++ {
					hsvs[i].H = uint8(time.Since(start).Seconds())
					hsvs[i].V = v
				}

				c.hal.Leds().HSV(0, hsvs...)
				c.hal.Leds().Show()
			}
		}
	}()

	c.hal.LCD().Write(c.errorMessage[0], c.errorMessage[1], hal.LCDRed)
	time.Sleep(time.Second)
	for {
		if c.hal.Button() {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
	return stateReady
}
