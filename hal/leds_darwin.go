package hal

import (
	"sync"
)

type leds struct {
	value bool
	lock  sync.Mutex
}

func newLeds() *leds {
	return &leds{}
}

func (l *leds) HSV(i int, values ...HSV) {

}
func (l *leds) RGB(i int, values ...RGB) {

}
func (l *leds) Show() {

}
func (l *leds) Clear() {

}
