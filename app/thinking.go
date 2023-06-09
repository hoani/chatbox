package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hoani/chatbox/hal"
	"github.com/hoani/chatbox/lcd"
	openai "github.com/sashabaranov/go-openai"
)

func (c *chatbox) doStateThinking() state {
	c.hal.LCD().Write(lcd.Pad("[Thinking]"), "", hal.LCDBlue)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		hsvs := []hal.HSV{}
		for i := 0; i < 12; i++ {
			hsvs = append(hsvs, hal.HSV{
				H: 0xa0,
				S: 0xFF,
				V: 20 + 20*uint8(i),
			})
		}
		for i := 0; i < 12; i++ {
			hsvs = append(hsvs, hal.HSV{
				H: 0xf0,
				S: 0xFF,
				V: 20 + 20*uint8(i),
			})
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Millisecond * 50):
				last := hsvs[23]
				for i := 23; i >= 0; i-- {
					hsvs[i] = hsvs[(i+23)%24]
				}
				hsvs[0] = last

				c.hal.Leds().HSV(0, hsvs...)
				c.hal.Leds().Show()
			}
		}
	}()

	var path string
	select {
	case path = <-c.recordingCh:
		if path == "" {
			c.hal.Debug("recording path is empty")
			return stateReady
		}
	case <-time.After(time.Second * 5):
		c.hal.Debug("timeout waiting for recording")
		return stateReady
	}

	transcription, err := c.openai.CreateTranscription(
		ctx,
		openai.AudioRequest{
			Model:    openai.Whisper1,
			FilePath: path,
		})
	if err != nil {
		c.hal.Debug(fmt.Sprintf("transcription error: %#v\n", err))
		return stateReady
	}

	c.hal.Debug(fmt.Sprintf("User: %s \n", transcription.Text))

	if next, ok := c.handleUserCommands(transcription.Text); ok {
		return next
	}

	// Moderate messages.
	if ok := c.moderator.Moderate(transcription.Text); !ok {
		c.errorMessage = "I can't talk about this"
		return stateError
	}

	c.chat.Messages = append(c.chat.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: transcription.Text,
	})

	c.refreshChatRequest()
	resp, err := c.openai.CreateChatCompletion(ctx, *c.chat)
	if err != nil {
		c.hal.Debug(fmt.Sprintf("chat error: %#v\n", err))
		c.chat.Messages = c.chat.Messages[:len(c.chat.Messages)-1] // Remove the last message.
		return stateReady
	}

	c.hal.Debug(resp.Choices[0].Message.Content)
	c.chat.Messages = append(c.chat.Messages, resp.Choices[0].Message)
	c.lastChat = time.Now() // Success, update last chat time.

	return stateTalking
}

func (c *chatbox) handleUserCommands(input string) (state, bool) {
	input = strings.TrimSuffix(strings.ToLower(input), ".")
	if input == "shutdown" || input == "shut down" {
		return stateShutdown, true
	}
	if strings.HasPrefix(input, "change personality to ") {
		c.personality = strings.TrimPrefix(input, "change personality to ")
		return stateChange, true
	}

	return stateThinking, false
}
