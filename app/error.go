package app

import (
	"strings"
	"time"

	"github.com/hoani/chatbox/hal"
)

func (c *chatbox) doStateError() state {
	start := time.Now()
	msg := strings.Repeat(" ", 16) + c.errorMessage + strings.Repeat(" ", 16)
	for {
		if c.hal.Button() && time.Since(start) > time.Second {
			break
		}
		time.Sleep(time.Millisecond * 10)
		start := int(time.Since(start)/(time.Millisecond*100)) % (len(msg) - 16)
		end := start + 16
		c.hal.LCD().Write("    [Error]    ", msg[start:end], hal.LCDRed)
	}
	return stateReady
}