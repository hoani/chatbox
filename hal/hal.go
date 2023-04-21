package hal

type LCDColor uint

const (
	LCDRed LCDColor = iota
	LCDGreen
	LCDAqua
	LCDBlue
)

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
	Write(line1, line2 string, color LCDColor)
}

type Hal interface {
	Button() bool
	Leds() Leds
	LCD() LCD
	Debug(string)
	Shutdown()
}

func NewHal() (Hal, error) {
	return newHal()
}
