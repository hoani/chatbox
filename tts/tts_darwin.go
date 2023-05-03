package tts

func getFlags(c Config) []string {
	flags := []string{}
	if c.Male {
		if c.AltVoice {
			flags = append(flags, "-v", "m7")
		} else {
			flags = append(flags, "-v", "en")
		}
	} else {
		if c.AltVoice {
			flags = append(flags, "-v", "f2")
		} else {
			flags = append(flags, "-v", "f1")
		}
	}
	return flags
}
