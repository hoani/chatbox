package hal

import "github.com/stianeikeland/go-rpio/v4"

const buttonPin = 5

type hal struct {
	button rpio.Pin
	leds *leds
}

func newHal() (*hal, error) {
	if err := rpio.Open(); err != nil {
		return nil, err
	}

	button := rpio.Pin(buttonPin)
	button.Input()
	button.PullUp()

	leds, err := newLeds()
	if err != nil {
		return nil, err
	}

	return &hal{
		button: button,
		leds: leds,
	}, nil
}

func (h *hal) Button() bool {
	return h.button.Read() == rpio.Low // Note: reverse polarity, High means unpressed.
}

func (h *hal) Leds() Leds {
	return h.leds
}