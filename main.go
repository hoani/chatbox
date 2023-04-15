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
