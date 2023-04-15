package app

import (
	"context"
	"fmt"
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

	translation, err := c.openai.CreateTranslation(
		context.Background(),
		openai.AudioRequest{
			Model:    openai.Whisper1,
			FilePath: path,
		})
	if err != nil {
		c.hal.Debug(fmt.Sprintf("translation error: %#v\n", err))
		return stateReady
	}

	c.hal.Debug(fmt.Sprintf("User: %s \n", translation.Text))

	c.chat.Messages = append(c.chat.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: translation.Text,
	})

	c.refreshChatRequest()
	resp, err := c.openai.CreateChatCompletion(context.Background(), *c.chat)
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
