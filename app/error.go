package app

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/hoani/chatbox/hal"
	"github.com/hoani/chatbox/lcd"
)

const errorLCDUpdateRate = time.Millisecond * 100

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

	msg := strings.Repeat(" ", 16) + c.errorMessage + strings.Repeat(" ", 16)
	index := 0
	lcdLast := time.Now()
	start := time.Now()

	for {
		if time.Since(lcdLast) >= errorLCDUpdateRate {
			lcdLast = lcdLast.Add(errorLCDUpdateRate)
			c.hal.LCD().Write(lcd.Pad("[error]"), msg[index:index+15], hal.LCDRed)
			index = (index + 1) % (len(msg) - 16)
		}

		if c.hal.Button() && time.Since(start) > time.Second {
			break
		}
		time.Sleep(time.Millisecond * 10)

	}
	return stateReady
}
