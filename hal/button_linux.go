package hal

import "github.com/stianeikeland/go-rpio/v4"

const buttonPin = 5

type hal struct {
	button rpio.Pin
}

func newHal() (*hal, error) {
	if err := rpio.Open(); err != nil {
		return nil, err
	}

	button := rpio.Pin(buttonPin)
	button.Input()
	button.PullUp()

	return &hal{
		button: button,
	}, nil
}

func (h *hal) Button() bool {
	return h.Read() == rpio.Low // Note: reverse polarity, High means unpressed.
}
