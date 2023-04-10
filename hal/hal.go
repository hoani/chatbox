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

type LCD interface {
	Write(line1, line2 string, color *RGB)
}

type Hal interface {
	Button() bool
	Leds() Leds
	LCD() LCD
}

func NewHal() (Hal, error) {
	return newHal()
}
