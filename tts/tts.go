package tts

import (
	"os/exec"
)

type Config struct {
	Male     bool
	AltVoice bool
}

func Speak(input string, c Config) error {
	args := append(getFlags(c), `"`+input+`"`, "-z")
	return exec.Command("espeak", args...).Run()
}
