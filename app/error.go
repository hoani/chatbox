package app

import (
	"time"

	"github.com/hoani/chatbox/hal"
)

func (c *chatbox) doStateError() state {
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
