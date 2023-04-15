package app

import (
	"strings"
	"time"

	"github.com/hoani/chatbox/hal"
)

func (c *chatbox) doStateError() state {
	start := time.Now()
	for {
		if c.hal.Button() && time.Since(start) > time.Second {
			break
		}
		time.Sleep(time.Millisecond * 10)
		start := int(time.Since(start)/(time.Millisecond*200)) % (len(msg) - 16)
		end := start + 16
		c.hal.LCD().Write(c.errorMessage[0], c.errorMessage[1], hal.LCDRed)
	}
	return stateReady
}
