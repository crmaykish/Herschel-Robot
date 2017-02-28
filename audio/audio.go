package audio

import "os/exec"

// Sound plays back a wav file. Starts as a goroutine to allow normal operation to continue
func Sound(file string) {
	// TODO: check that file exists and that sound actually plays
	go exec.Command("aplay", "./sound/"+file+".wav", "-q").Run()
}
