package hal

import (
	"fmt"
	"strings"
	"os"
	"io"
	"errors"

	"github.com/tarm/goserial"
)

const numLeds = 24

type leds struct {
	s io.ReadWriteCloser
}

func newLeds() (*leds, error) {
	port, err := findPort()
	if err != nil {
		return nil, err
	}

	c := &serial.Config{Name: port, Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}

	// Initialize the LEDs to 24, clear and show.
	_, err = s.Write([]byte(fmt.Sprintf("I%02x\nC\nS\n", numLeds)))
	if err != nil {
			return nil, err
	}
	return &leds{
		s: s,
	}, nil
}

func findPort() (string, error) {
	contents, _ := os.ReadDir("/dev")

	// Look for what is mostly likely the Arduino device
	for _, f := range contents {
		if strings.Contains(f.Name(), "ttyACM") {
			return "/dev/" + f.Name(), nil
		}
	}

	// Have not been able to find a USB device that 'looks'
	// like an Arduino.
	return "", errors.New("cannot find USB device")
}


func (l *leds) HSV(i int, values ...HSV){
	if len(values) == 0 {
		return
	}
	msg := fmt.Sprintf("H%02x", i)
	for _, v := range values {
		msg += fmt.Sprintf("%02x%02x%02x", v.H, v.S, v.V)
	}
	msg += "\n"
	_, err := l.s.Write([]byte(msg))
	if err != nil {
		fmt.Println("error writing to serial " + err.Error())
	}
}
func (l *leds) RGB(i int, values ...RGB){
	if len(values) == 0 {
		return
	}
	msg := fmt.Sprintf("R%02x", i)
	for _, v := range values {
		msg += fmt.Sprintf("%02x%02x%02x", v.R, v.G, v.B)
	}
	msg += "\n"
	_, err := l.s.Write([]byte(msg))
	if err != nil {
		fmt.Println("error writing to serial " + err.Error())
	}
}
func (l *leds) Show(){
	_, err := l.s.Write([]byte("S\n"))
	if err != nil {
		fmt.Println("error writing to serial " + err.Error())
	}
}
func (l *leds) Clear(){
	_, err := l.s.Write([]byte("C\n"))
	if err != nil {
		fmt.Println("error writing to serial " + err.Error())
	}
}


