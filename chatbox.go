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

	"github.com/hoani/chatbox/3rdparty/faiface/beep/wav"
	"github.com/hoani/chatbox/hal"
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
	openai *openai.Client
	hal    hal.Hal
	wd     string
	state  state
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
		openai: c,
		hal:    h,
		wd:     wd,
		state:  stateReady,
	}, nil
}

func (c *chatbox) getRecording() string {
	m, err := toot.NewDefaultMicrophone()
	if err != nil {
		panic(err)
	}
	defer m.Close()

	// h.Debug(fmt.Sprintf("%#v\n", m.DeviceInfo()))

	path := filepath.Join(c.wd, "test.wav")
	f, err := os.Create(path)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err = wav.Encode(f, m, m.Format())
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := m.Start(ctx); err != nil {
		panic(err)
	}

	for {
		if !c.hal.Button() {
			break
		}
		time.Sleep(time.Millisecond)
	}
	time.Sleep(time.Millisecond * 100)
	m.Close()
	wg.Wait()

	resp, err := c.openai.CreateTranslation(
		context.Background(),
		openai.AudioRequest{
			Model:    openai.Whisper1,
			FilePath: path,
		})
	if err != nil {
		panic(err)
	}
	return resp.Text
}

func (c *chatbox) run() error {
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: "respond as an exaggerated Jim Carrey whose soul has been trapped inside a raspberry pi. " +
					"When the user calls you by an incorrect name, respond as if they said your name correctly. ",
			},
		},
		Temperature: 1.0,
	}

	hsvs := []hal.HSV{}
	for i := 0; i < 24; i++ {
		hsvs = append(hsvs, hal.HSV{
			H: uint8(i) * 10,
			S: 0xFF,
			V: 0x50,
		})
	}

	for {
		for {
			if c.hal.Button() {
				break
			}

			for i := range hsvs {
				hsvs[i].H += 1
			}

			c.hal.Leds().HSV(0, hsvs...)
			c.hal.Leds().Show()

			time.Sleep(10 * time.Millisecond)

		}

		c.hal.LCD().Write("Listening...", "release to stop", &hal.RGB{R: 200, G: 205, B: 0})

		input := c.getRecording()
		c.hal.LCD().Write("Thinking...", "", &hal.RGB{R: 0, G: 205, B: 0})
		req.Messages = append(req.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: input,
		})

		c.hal.Debug(fmt.Sprintf("User: %s \n\n", input))
		resp, err := c.openai.CreateChatCompletion(context.Background(), req)
		if err != nil {
			panic(err)
		}
		c.hal.Debug(fmt.Sprintf("%s\n\n", resp.Choices[0].Message.Content))
		c.hal.LCD().Write("Talking...", "", &hal.RGB{R: 0, G: 205, B: 100})

		cmd := exec.Command("espeak", `"`+resp.Choices[0].Message.Content+`"`)
		cmd.Run()

		c.hal.LCD().Write("Press to start", "", &hal.RGB{R: 0, G: 205, B: 200})
		req.Messages = append(req.Messages, resp.Choices[0].Message)

	}
}
