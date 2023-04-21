package app

import (
	"fmt"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

const systemMsgBase = "Respond as an exaggerated %s whose soul is trapped inside a raspberry Pi. " +
	"When possible keep responses to less than three sentences. "

func (c *chatbox) newChatRequest() *openai.ChatCompletionRequest {
	return c.newCustomChatRequest("Jim Carrey")
}

func (c *chatbox) systemMessage() string {
	return fmt.Sprintf(systemMsgBase, c.personality)
}

func (c *chatbox) newCustomChatRequest(personality string) *openai.ChatCompletionRequest {
	c.personality = personality
	return &openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: c.systemMessage(),
			},
		},
		Temperature: 1.0,
	}
}

func (c *chatbox) refreshChatRequest() {
	c.chat.Messages[0].Content = c.systemMessage()
	if time.Since(c.lastChat) < time.Hour*24*365 {
		c.chat.Messages[0].Content += "the last time we spoke was " + c.lastChat.Format(time.Stamp) + "."
	}
	c.chat.Messages[0].Content += " the current time is " + time.Now().Format(time.Stamp) + "."
}
