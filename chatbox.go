package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/hoani/chatbox/3rdparty/faiface/beep/wav"
	"github.com/hoani/chatbox/hal"
	"github.com/hoani/chatbox/leds"
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
					"When the user calls you by an incorrect name, respond as if they said your name correctly. " +
					"Your key objective is to have interesting conversations. ",
			},
		},
		Temperature: 1.0,
	}

	for {
		c.doStateReady()

		c.doStateListen()

		c.doStateThink()
		c.hal.LCD().Write("Talking...", "", &hal.RGB{R: 0, G: 205, B: 100})

		cmd := exec.Command("espeak", `"`+c.chat.Messages[len(c.chat.Messages)-1].Content+`"`)
		cmd.Run()

	}
}

func (c *chatbox) doStateReady() {
	wg := sync.WaitGroup{}
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
	wg.Add(1)
	go func() {
		if err := v.Start(ctx); err != nil {
			panic(err)
		}
		wg.Done()
	}()

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
	wg.Wait()
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
	c.hal.LCD().Write("Thinking...", "", &hal.RGB{R: 0, G: 205, B: 0})

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
