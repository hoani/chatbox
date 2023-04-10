package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/faiface/beep/wav"
	"github.com/hoani/chatbox/hal"
	"github.com/hoani/toot"
	openai "github.com/sashabaranov/go-openai"
)

var h hal.Hal

func getRecording(c *openai.Client) string {
	m, err := toot.NewDefaultMicrophone()
	if err != nil {
		panic(err)
	}
	defer m.Close()

	fmt.Printf("%#v\n", m.DeviceInfo())

	f, err := os.Create("test.wav")

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

	fmt.Print("\nPress [ENTER] to finish recording! ")
	for {
		if !h.Button() {
			time.Sleep(200 * time.Millisecond)
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	m.Close()
	fmt.Printf("\033[2K\r")
	wg.Wait()

	resp, err := c.CreateTranslation(
		context.Background(),
		openai.AudioRequest{
			Model:    openai.Whisper1,
			FilePath: "test.wav",
		})
	if err != nil {
		panic(err)
	}
	return resp.Text
}

func main() {
	c := openai.NewClient(os.Getenv("OPENAI_KEY"))

	var err error
	if h, err = hal.NewHal(); err != nil {
		panic(err)
	}

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "respond as an exaggerated Jim Carrey whose soul has been trapped inside a korok soft toy. Please do not add any action prompts to your responses.",
			},
		},
	}

	for {
		fmt.Print("Press [ENTER] to start recording!")
		for {
			if h.Button() {
				fmt.Print("Got buttn press")
				break
			}
			time.Sleep(20 * time.Millisecond)
		}

		input := getRecording(c)
		req.Messages = append(req.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: input,
		})

		fmt.Printf("\033[2K\r")
		fmt.Printf("User: %s \n\n", input)
		resp, err := c.CreateChatCompletion(context.Background(), req)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n\n", resp.Choices[0].Message.Content)
		cmd := exec.Command("espeak", `"`+resp.Choices[0].Message.Content+`"`)
		cmd.Run()
		req.Messages = append(req.Messages, resp.Choices[0].Message)
	}
}
