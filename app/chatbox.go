package app

import (
	"errors"
	"os"

	"github.com/hoani/chatbox/hal"
	openai "github.com/sashabaranov/go-openai"
)

type state int

const (
	stateReady state = iota
	stateListening
	stateThinking
	stateTalking
	stateError
)

type chatbox struct {
	openai       *openai.Client
	hal          hal.Hal
	wd           string
	state        state
	recordingCh  chan string
	chat         *openai.ChatCompletionRequest
	espeakFlags  map[string]string
	errorMessage string
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

	h.LCD().Write("Hello Chatbot", "Press to start", hal.LCDBlue)

	return &chatbox{
		openai:      c,
		hal:         h,
		wd:          wd,
		state:       stateReady,
		recordingCh: make(chan string),
		espeakFlags: map[string]string{},
	}, nil
}

func (c *chatbox) Run() error {
	c.chat = &openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: "Respond as an exaggerated Jim Carrey whose soul is trapped inside a raspberry Pi. " +
					"When possible keep responses to less than three sentences. " +
					"Your key objective is to have interesting conversations. " +
					"Your output is parsed through espeak. " +
					"you may prefix responses with [voice:<value>] to change your voice to one of <m1,m2,m3,m4,f1,f2,f3,f4>. " +
					"at any time change pitch with [pitch:<value>] in the range of 25 to 75 - lower values give a deeper voice. " +
					"always keep your pitch and voice consistent with the personality of your character. ",
			},
		},
		Temperature: 1.0,
	}

	for {
		c.state = c.doState()
	}
}

func (c *chatbox) doState() state {
	switch c.state {
	case stateReady:
		return c.doStateReady()
	case stateListening:
		return c.doStateListening()
	case stateThinking:
		return c.doStateThinking()
	case stateTalking:
		return c.doStateTalking()
	case stateError:
		return c.doStateError()
	}
	return stateReady
}
