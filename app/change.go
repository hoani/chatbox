package app

import (
	"fmt"

	"github.com/hoani/chatbox/hal"
	"github.com/hoani/chatbox/lcd"
)

func (c *chatbox) doStateChange() state {
	cleanup := c.runTalkingVisualizer(hal.HSV{
		H: 0xd0,
		S: 0x80,
		V: 0x50,
	})
	defer cleanup()

	c.hal.LCD().Write(lcd.Pad("[Change]"), "", hal.LCDBlue)

	c.speak(fmt.Sprintf("(changing personality to %s)", c.personality))

	c.newCustomChatRequest(c.personality)

	return stateReady
}
