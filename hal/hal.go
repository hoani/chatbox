package hal

type Hal interface {
	Button() bool
}

func NewHal() (Hal, error) {
	return newHal()
}
