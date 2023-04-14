package main

import (
	"os"
)

func main() {

	key := os.Getenv("OPENAI_KEY")
	cb, err := NewChatBox(key)
	if err != nil {
		panic(err)
	}

	if err := cb.run(); err != nil {
		panic(err)
	}
}
