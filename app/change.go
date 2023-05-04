package app

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hoani/chatbox/hal"
	"github.com/hoani/chatbox/lcd"
	openai "github.com/sashabaranov/go-openai"
)

func (c *chatbox) doStateChange() state {
	cleanup := c.runTalkingVisualizer(hal.HSV{
		H: 0xd0,
		S: 0x80,
		V: 0x50,
	})
	defer cleanup()

	c.hal.LCD().Write(lcd.Pad("[Change]"), "", hal.LCDBlue)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.getGender()
	}()

	c.processDirective(fmt.Sprintf("(changing personality to %s)", c.personality))
	c.chat = c.newCustomChatRequest(c.personality)

	wg.Wait()

	return stateReady
}

func (c *chatbox) getGender() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := fmt.Sprintf("Respond with only 'Male' or 'Female'. What gender is %s?", c.personality)
	res, err := c.openai.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: query,
				},
			},
			Temperature: 0.2,
		},
	)
	if err != nil || len(res.Choices) == 0 {
		return
	}
	c.ttsCfg.Male = !strings.Contains(strings.ToUpper(res.Choices[0].Message.Content), "FEMALE")
}
