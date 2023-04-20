package app

import (
	"context"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hoani/chatbox/hal"
	"github.com/hoani/chatbox/lcd"
	"github.com/hoani/chatbox/leds"
	"github.com/hoani/chatbox/strutil"
	openai "github.com/sashabaranov/go-openai"
)

func (c *chatbox) doStateTalking() state {
	message := c.chat.Messages[len(c.chat.Messages)-1]
	if message.Role != openai.ChatMessageRoleAssistant {
		return stateReady // Oops, this isn't a response, best get out of here.
	}
	content := message.Content

	// Ideally, we would use the audio out rather than microphone... but this works well anyway.
	v := leds.NewVisualizer()

	ctx, cancel := context.WithCancel(context.Background())
	go v.Start(ctx)
	defer v.Wait()
	defer cancel()

	hsvChan := make(chan hal.HSV)

	go func() {
		baseHsv := hal.HSV{
			H: 0x80,
			S: 0x00,
			V: 0x50,
		}

		hsvs := make([]hal.HSV, 24)

		for {
			select {
			case <-ctx.Done():
				return
			case baseHsv = <-hsvChan:
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

	c.hal.LCD().Write(lcd.Pad("[Talking]"), "", hal.LCDBlue)
	directives := strutil.SplitBrackets(content)
	for _, directive := range directives {
		c.processDirective(directive)
	}

	return stateReady
}

func (c *chatbox) processDirective(d string) {
	if strings.HasPrefix(d, "[") {
		c.processCommandBlock(d)
		return
	}
	if strings.HasPrefix(d, "(") || strings.HasPrefix("*") {
		c.espeakFlags["-v"] = "m7"
	} else {
		c.espeakFlags["-v"] = "en"
	}
	c.processSpeech(d)
}

func (c *chatbox) processCommandBlock(block string) {
	block = strings.TrimPrefix(block, "[")
	block = strings.TrimSuffix(block, "]")
	cmds := strings.Split(block, " ")
	for _, cmd := range cmds {
		parts := strings.Split(cmd, ":")
		if len(parts) != 2 {
			continue
		}
		switch parts[0] {
		case "pitch":
			if isValidPitch(parts[1]) {
				c.espeakFlags["-p"] = parts[1]
			}
		}
	}
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
		wpm := 175
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
	args := []string{`"` + sentence + `"`, "-z"}
	for flag, arg := range c.espeakFlags {
		args = append(args, flag, arg)
	}
	// c.hal.Debug(strings.Join(args, " "))
	exec.Command("espeak", args...).Run()
}

func isValidPitch(in string) bool {
	v, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		return false
	}
	if v < 0 || v > 99 {
		return false
	}
	return true
}
