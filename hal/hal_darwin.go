package hal

type hal struct {
	button *button
	leds   *leds
	lcd    *lcd
}

func newHal() (*hal, error) {
	h := &hal{
		button: &button{},
		lcd:    &lcd{},
		leds:   &leds{},
	}
	h.button.start()
	return h, nil
}

func (h *hal) Leds() Leds {
	return h.leds
}

func (h *hal) LCD() LCD {
	return h.lcd
}

func (h *hal) Button() bool {
	return h.button.get()
}
