package app

import (
	"errors"
	"os"
	"time"

	"github.com/hoani/chatbox/hal"
	openai "github.com/sashabaranov/go-openai"
)

type state int

const buttonDebounce = time.Millisecond * 50

const (
	stateReady state = iota
	stateListening
	stateThinking
	stateTalking
	stateShutdown
	stateChange
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
	lastChat     time.Time
	personality  string
}

func NewChatBox(key string) (*chatbox, error) {
	if key == "" {
		return nil, errors.New("missing openai key")
	}
	c := openai.NewClient(key)

	h, err := hal.NewHal()
	if err != nil {
		return nil, err
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
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
		espeakFlags: map[string]string{
			"-v": "m7",
		},
	}, nil
}

func (c *chatbox) Run() error {
	c.chat = c.newChatRequest()

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
	case stateShutdown:
		return c.doStateShutdown()
	case stateChange:
		return c.doStateChange()
	}
	return stateReady
}
