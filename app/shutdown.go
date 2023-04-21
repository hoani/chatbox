package app

import (
	"github.com/hoani/chatbox/hal"
	"github.com/hoani/chatbox/lcd"
)

func (c *chatbox) doStateShutdown() state {
	cleanup := c.runTalkingVisualizer(hal.HSV{
		H: 0xd0,
		S: 0x80,
		V: 0x50,
	})
	defer cleanup()

	c.hal.LCD().Write(lcd.Pad("[Shutdown]"), "", hal.LCDBlue)

	c.speak("(shutting down)")

	c.hal.Shutdown()

	return stateReady
}
