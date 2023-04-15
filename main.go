package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hoani/chatbox/app"
	"github.com/hoani/chatbox/hal"
	"github.com/hoani/chatbox/strutil"
)

func main() {
	key := os.Getenv("OPENAI_KEY")
	cb, err := app.NewChatBox(key)
	if err != nil {
		handleError(err)
	}

	if err := cb.Run(); err != nil {
		handleError(err)
	}
}

// Attempts to report error by various means.
func handleError(err error) {
	fmt.Println("program error: " + err.Error())

	// Cool down for a while. We don't want to smash the processor with restart attempts.
	timer := time.NewTimer(time.Hour)
	h, e := hal.NewHal()
	if e != nil {
		<-timer.C
		os.Exit(1)
	}

	rgbs := []hal.RGB{}
	for i := 0; i < 24; i++ {
		rgbs = append(rgbs, hal.RGB{R: 0xFF, G: 0, B: 0})
	}
	h.Leds().RGB(0, rgbs...)
	h.Leds().Show()

	msgs := strutil.SplitWidth(err.Error(), 16)
	msgs = append(msgs, "")
	if len(msgs)%2 != 0 {
		msgs = append(msgs, "")
	}
	index := 0

	for {
		select {
		case <-timer.C:
			os.Exit(1)
		case <-time.After(time.Second):
			h.LCD().Write(msgs[index%len(msgs)], msgs[(index+1)%len(msgs)], hal.LCDBlue)
			index += 1
		}
	}
}
