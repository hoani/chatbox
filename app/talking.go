package app

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/hoani/chatbox/hal"
	"github.com/hoani/chatbox/lcd"
	"github.com/hoani/chatbox/leds"
	"github.com/hoani/chatbox/strutil"
	"github.com/hoani/chatbox/tts"
	openai "github.com/sashabaranov/go-openai"
)

func (c *chatbox) doStateTalking() state {
	message := c.chat.Messages[len(c.chat.Messages)-1]
	if message.Role != openai.ChatMessageRoleAssistant {
		return stateReady // Oops, this isn't a response, best get out of here.
	}
	content := message.Content

	cleanup := c.runTalkingVisualizer(hal.HSV{
		H: 0x80,
		S: 0x00,
		V: 0x50,
	})
	defer cleanup()

	c.hal.LCD().Write(lcd.Pad("[Talking]"), "", hal.LCDBlue)
	directives := strutil.SplitBrackets(content)
	for _, directive := range directives {
		c.processDirective(directive)
	}

	return stateReady
}

func (c *chatbox) processDirective(d string) {
	if strings.HasPrefix(d, "[") {
		return
	}
	if strings.HasPrefix(d, "(") || strings.HasPrefix(d, "*") {
		c.ttsCfg.AltVoice = true
	} else {
		c.ttsCfg.AltVoice = false
	}
	c.processSpeech(d)
}

func (c *chatbox) processSpeech(in string) {
	sentences := strutil.SplitSentences(in)
	for _, sentence := range sentences {
		parts := strutil.SplitWidth(sentence, 16)
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.speak(sentence)
		}()
		wpm := 155
		adjustment := 0.6 // espeak is a bit faster than the wpm would suggest.
		mspw := int(adjustment * (60.0 * 1000.0) / float64(wpm))
		for _, part := range parts {
			part = strings.TrimSpace(part)
			words := strings.Count(part, " ") + 1
			padding := (16 - len(part)) / 2
			part = strings.Repeat(" ", padding) + part
			c.hal.LCD().Write(lcd.Pad("[Talking]"), part, hal.LCDAqua)
			time.Sleep(time.Millisecond * time.Duration(mspw*words))
		}
		wg.Wait()
	}
}

func (c *chatbox) speak(sentence string) {
	tts.Speak(sentence, tts.Config{Male: true})
}

func (c *chatbox) runTalkingVisualizer(baseHsv hal.HSV) (cleanup func()) {
	// Ideally, we would use the audio out rather than microphone... but this works well anyway.
	v := leds.NewVisualizer()

	ctx, cancel := context.WithCancel(context.Background())
	go v.Start(ctx)

	go func() {
		hsvs := make([]hal.HSV, 24)

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Millisecond * 50):
				channels := v.Channels()
				for i := range hsvs {
					hsvs[i] = baseHsv
					hsvs[i].V = 0x50 + uint8(channels[i%leds.NChannels])
				}

				c.hal.Leds().HSV(0, hsvs...)
				c.hal.Leds().Show()
			}
		}
	}()

	return func() {
		cancel()
		v.Wait()
	}
}
