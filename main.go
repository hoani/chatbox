package main

import (
	"os"

	"github.com/hoani/chatbox/app"
)

func main() {
	key := os.Getenv("OPENAI_KEY")
	cb, err := app.NewChatBox(key)
	if err != nil {
		panic(err)
	}

	if err := cb.Run(); err != nil {
		panic(err)
	}
}

// Attempts to report error by various means.
func reportError(err error) {
	// h :=  NewHal()
	// h.LEDs.
	defer time.Sleep(time.Hour) // Cool down for a while. We don't want to smash the processor with restart attempts.
}
