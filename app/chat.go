package app

import (
	"time"

	openai "github.com/sashabaranov/go-openai"
)

const systemMsgBase = "Respond as an exaggerated Jim Carrey whose soul is trapped inside a raspberry Pi. " +
	"When possible keep responses to less than three sentences. " +
	"Your key objective is to have interesting conversations. " +
	"You may change pitch by prefixing sentences with [pitch:<value>] in the range of 25 to 75 - lower values give a deeper voice. " +
	"change pitch to express emotion"

func (c *chatbox) newChatRequest() *openai.ChatCompletionRequest {
	return &openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemMsgBase,
			},
		},
		Temperature: 1.0,
	}
}

func (c *chatbox) refreshChatRequest() {
	c.chat.Messages[0].Content = systemMsgBase
	if time.Since(c.lastChat) < time.Hour*24*365 {
		c.chat.Messages[0].Content += "the last time we spoke was " + c.lastChat.Format(time.Stamp) + "."
	}
	c.chat.Messages[0].Content += " the current time is " + time.Now().Format(time.Stamp) + "."
}
