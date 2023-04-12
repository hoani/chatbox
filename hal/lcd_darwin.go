package hal

type lcd struct{}

func newLCD() *lcd {
	return &lcd{}
}

func (l *lcd) Write(line1, line2 string, color *RGB) {

}
