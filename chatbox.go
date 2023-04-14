package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/hoani/chatbox/3rdparty/faiface/beep/wav"
	"github.com/hoani/chatbox/hal"
	"github.com/hoani/chatbox/leds"
	"github.com/hoani/chatbox/strutil"
	"github.com/hoani/toot"
	openai "github.com/sashabaranov/go-openai"
)

type state int

const (
	stateReady state = iota
	stateListening
	stateThinking
	stateTalking
)

type chatbox struct {
	openai      *openai.Client
	hal         hal.Hal
	wd          string
	state       state
	recordingCh chan string
	chat        *openai.ChatCompletionRequest
}

func NewChatBox(key string) (*chatbox, error) {
	if key == "" {
		return nil, errors.New("missing openai key")
	}
	c := openai.NewClient(key)

	h, err := hal.NewHal()
	if err != nil {
		panic(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	hsvs := []hal.HSV{}
	for i := 0; i < 24; i++ {
		hsvs = append(hsvs, hal.HSV{
			H: uint8(i) * 10,
			S: 0xFF,
			V: 0x50,
		})
	}

	h.Leds().HSV(0, hsvs...)
	h.Leds().Show()

	h.LCD().Write("Hello Chatbot", "Press to start", &hal.RGB{R: 100, G: 105, B: 200})

	return &chatbox{
		openai:      c,
		hal:         h,
		wd:          wd,
		state:       stateReady,
		recordingCh: make(chan string),
	}, nil
}

func (c *chatbox) run() error {
	c.chat = &openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: "respond as an exaggerated Jim Carrey whose soul has been trapped inside a raspberry pi. " +
					"When possible keep responses to less than three sentences. " +
					"Your key objective is to have interesting conversations. ",
			},
		},
		Temperature: 1.0,
	}

	for {
		c.doStateReady()

		c.doStateListen()

		c.doStateThink()

		c.doStateTalking()

	}
}

func (c *chatbox) doStateReady() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c.hal.LCD().Write("Press to start", "", &hal.RGB{R: 0, G: 205, B: 200})

	hsvs := []hal.HSV{}
	for i := 0; i < 24; i++ {
		hsvs = append(hsvs, hal.HSV{
			H: uint8(i) * 10,
			S: 0xFF,
			V: 0x50,
		})
	}

	v := leds.NewVisualizer()
	go func() {
		if err := v.Start(ctx); err != nil {
			panic(err)
		}
	}()
	defer v.Wait()

	for {
		if c.hal.Button() {
			cancel()
			break
		}

		time.Sleep(20 * time.Millisecond)

		channels := v.Channels()

		for i := range hsvs {
			hsvs[i].H += 1
			j := i
			if j >= leds.NChannels {
				j = leds.NChannels - (1 + i - leds.NChannels)
			}
			v := channels[j]
			if v > float64(0xa0) {
				v = float64(0xa0)
			}
			hsvs[i].V = 0x40 + uint8(v)
			// c.hal.Debug(fmt.Sprintf("%#v\n", channels))
		}

		c.hal.Leds().HSV(0, hsvs...)
		c.hal.Leds().Show()
	}
}

func (c *chatbox) doStateListen() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c.hal.LCD().Write("Listening...", "release to stop", &hal.RGB{R: 200, G: 205, B: 0})

	path := filepath.Join(c.wd, "test.wav")
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	m, err := toot.NewDefaultMicrophone()
	if err != nil {
		panic(err)
	}
	defer func() {
		go func() {
			time.Sleep(time.Second) // Delay a bit before stopping.
			m.Close()
		}()
	}()

	v := leds.NewVisualizer(
		leds.WithSource(&leds.Source{
			Streamer:   m,
			SampleRate: m.Format().SampleRate,
		}),
		leds.WithSink(func(s beep.Streamer) {
			go func() {
				err = wav.Encode(f, s, m.Format())
				if err != nil {
					c.hal.Debug(fmt.Sprintf("error encoding wav: %v", err))
					path = ""
				}
				f.Close()
				c.hal.Debug(path)
				c.recordingCh <- path
			}()
		}),
	)
	go v.Start(ctx)

	// h.Debug(fmt.Sprintf("%#v\n", m.DeviceInfo()))

	if err := m.Start(ctx); err != nil {
		panic(err)
	}

	hsvs := []hal.HSV{}
	for i := 0; i < 24; i++ {
		hsvs = append(hsvs, hal.HSV{
			H: 0x60,
			S: 0xFF,
			V: 0x50,
		})
	}

	for {
		if !c.hal.Button() {
			break
		}
		time.Sleep(time.Millisecond)

		channels := v.Channels()

		for i := range hsvs {
			j := i
			if j >= leds.NChannels {
				j = leds.NChannels - (1 + i - leds.NChannels)
			}
			v := channels[j]
			if v > float64(0xa0) {
				v = float64(0xa0)
			}
			hsvs[i].V = 0x40 + uint8(v)
		}

		c.hal.Leds().HSV(0, hsvs...)
		c.hal.Leds().Show()

	}
}

func (c *chatbox) doStateThink() {
	c.hal.LCD().Write("Thinking...", "", &hal.RGB{R: 0, G: 205, B: 200})

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
			return
		}
	case <-time.After(time.Second * 5):
		c.hal.Debug("timeout waiting for recording")
		return
	}

	translation, err := c.openai.CreateTranslation(
		context.Background(),
		openai.AudioRequest{
			Model:    openai.Whisper1,
			FilePath: path,
		})
	if err != nil {
		c.hal.Debug(fmt.Sprintf("translation error: %#v\n", err))
		return
	}

	c.hal.Debug(fmt.Sprintf("User: %s \n", translation.Text))

	c.chat.Messages = append(c.chat.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: translation.Text,
	})

	resp, err := c.openai.CreateChatCompletion(context.Background(), *c.chat)
	if err != nil {
		c.hal.Debug(fmt.Sprintf("chat error: %#v\n", err))
		return
	}
	c.hal.Debug(fmt.Sprintf("%s", resp.Choices[0].Message.Content))
	c.chat.Messages = append(c.chat.Messages, resp.Choices[0].Message)
}

func (c *chatbox) doStateTalking() {
	message := c.chat.Messages[len(c.chat.Messages)-1]
	if message.Role != openai.ChatMessageRoleAssistant {
		return // Oops, this isn't a response, best get out of here.
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

	c.hal.LCD().Write("   [Talking]   ", "", &hal.RGB{R: 0, G: 205, B: 100})
	directives := strutil.SplitBrackets(content)
	for _, directive := range directives {
		sentences := strutil.SplitSentences(directive)
		for _, sentence := range sentences {
			parts := strutil.SplitWidth(sentence, 16)
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				exec.Command("espeak", "-z", `"`+sentence+`"`).Run()
			}()
			wpm := 175
			adjustment := 0.8 // espeak is a bit faster than the wpm would suggest.
			mspw := int(adjustment * (60.0 * 1000.0) / float64(wpm))
			for _, part := range parts {
				part = strings.TrimSpace(part)
				words := strings.Count(part, " ") + 1
				padding := (16 - len(part)) / 2
				part = strings.Repeat(" ", padding) + part + strings.Repeat(" ", padding)
				c.hal.LCD().Write("   [Talking]   ", part, &hal.RGB{R: 0, G: 205, B: 100})
				time.Sleep(time.Millisecond * time.Duration(mspw*words))
			}
			wg.Wait()
		}
	}
}
