package main

import (
	"context"
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

var h hal.Hal

var wd string

func getRecording(c *openai.Client) string {
	m, err := toot.NewDefaultMicrophone()
	if err != nil {
		panic(err)
	}
	defer m.Close()

	// h.Debug(fmt.Sprintf("%#v\n", m.DeviceInfo()))

	path := filepath.Join(wd, "test.wav")
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
		if !h.Button() {
			time.Sleep(200 * time.Millisecond)
			break
		}
		time.Sleep(time.Millisecond)
	}
	m.Close()
	wg.Wait()

	resp, err := c.CreateTranslation(
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

func main() {

	key := os.Getenv("OPENAI_KEY")
	if key == "" {
		panic("Please set envvar OPENAI_KEY")
	}
	c := openai.NewClient(os.Getenv("OPENAI_KEY"))

	var err error
	if h, err = hal.NewHal(); err != nil {
		panic(err)
	}

	wd, err = os.Getwd()
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

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: "respond as an exaggerated Jim Carrey whose soul has been trapped inside a raspberry pi. " +
					"The raspberry pi is encased in a speaker with an LED ring and LCD display. " +
					"When the user calls you by an incorrect, but similar name, respond as if they said your name correctly. " +
					"Please do not add any action prompts to your responses." +
					"If the user has not introduced themselves, ask thier name. ",
			},
		},
		Temperature: 0.8,
	}

	for {
		for {
			if h.Button() {
				break
			}
			for i := range hsvs {
				hsvs[i].H += 1
			}

			h.Leds().HSV(0, hsvs...)
			h.Leds().Show()

			time.Sleep(10 * time.Millisecond)

		}

		h.LCD().Write("Listening...", "release to stop", &hal.RGB{200, 205, 0})

		input := getRecording(c)
		h.LCD().Write("Thinking...", "", &hal.RGB{0, 205, 0})
		req.Messages = append(req.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: input,
		})

		h.Debug(fmt.Sprintf("User: %s \n\n", input))
		resp, err := c.CreateChatCompletion(context.Background(), req)
		if err != nil {
			panic(err)
		}
		h.Debug(fmt.Sprintf("%s\n\n", resp.Choices[0].Message.Content))
		h.LCD().Write("Talking...", "", &hal.RGB{0, 205, 100})

		cmd := exec.Command("espeak", `"`+resp.Choices[0].Message.Content+`"`)
		cmd.Run()

		h.LCD().Write("Press to start", "", &hal.RGB{0, 205, 200})
		req.Messages = append(req.Messages, resp.Choices[0].Message)

	}
}
