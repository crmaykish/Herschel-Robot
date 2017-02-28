package main

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	xbox "github.com/crmaykish/Xboxdrv-Golang"
	"github.com/crmaykish/herschel/drive"
)

func mapRange(x, inMin, inMax, outMin, outMax int) int {
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

// Stop motors and disconnect serial connection
func stop() {
	if !drive.Connected {
		drive.Connect()
	}
	drive.Stop()
	drive.Disconnect()
}

func main() {
	// Watch for an OS interupt and trigger a cleanup
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		stop()
		os.Exit(1)
	}()

	fmt.Println("Starting Herschel control program...")

	xbox.Connect()

	go xbox.Control()

	drive.Connect()

	fmt.Println("Start control loop...")

	for {

		if math.Abs(float64(xbox.Xbox.LeftStick.Y)) < 1000 && math.Abs(float64(xbox.Xbox.RightStick.Y)) < 1000 {
			drive.Stop()
		} else {
			drive.Drive(mapRange(xbox.Xbox.LeftStick.Y, -32768, 32767, -255, 255), mapRange(xbox.Xbox.RightStick.Y, -32768, 32767, -255, 255))
		}

		// TODO: this should be a timer check to maintain a maximum of 60 Hz
		time.Sleep(time.Second / 60)
	}
}
