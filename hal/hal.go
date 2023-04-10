package hal

type HSV struct {
	H, S, V uint8
}

type RGB struct {
	R, G, B uint8
}

type Leds interface {
	HSV(i int, values ...HSV)
	RGB(i int, values ...RGB)
	Show()
	Clear()
}

type Hal interface {
	Button() bool
	Leds() Leds
}

func NewHal() (Hal, error) {
	return newHal()
}
