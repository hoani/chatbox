package hal

type hal struct {
	button *button
	leds   *leds
}

func newHal() (*hal, error) {
	h := &hal{
		button: &button{},
	}
	h.button.start()
	return h, nil
}

func (h *hal) Leds() Leds {
	return h.leds
}

func (h *hal) Button() bool {
	return h.button.get()
}
