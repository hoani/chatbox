package hal

import (
	"fmt"
	"sync"
)

type hal struct {
	button     bool
	buttonLock sync.Mutex
}

func newHal() (*hal, error) {
	h := &hal{
		button: false,
	}
	go func() {
		for {
			fmt.Scanln()
			h.buttonLock.Lock()
			h.button = !h.button
			h.buttonLock.Unlock()
		}
	}()
	return h, nil
}

func (h *hal) Button() bool {
	h.buttonLock.Lock()
	defer h.buttonLock.Unlock()
	return h.button
}
