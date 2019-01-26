package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crmaykish/herschel/pkg/audio"
	"github.com/crmaykish/herschel/pkg/drive"
	"github.com/crmaykish/herschel/pkg/lidar"
	xbox "github.com/crmaykish/xboxdrv-go"
)

// LoopRate is the rate to run the control loop at (in Hertz)
const LoopRate = 60

// MinJoystickSpeed before driving will start
const MinJoystickSpeed = 1000

// Map x from one range to another
func mapRange(x, inMin, inMax, outMin, outMax int) int {
	// TODO: check that x is within range first
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

// x is in range (inclusive)
func between(x, min, max int) bool {
	return x <= max && x >= min
}

// Stop motors and disconnect serial connection
func stop() {
	if !drive.Connected {
		drive.Connect()
	}
	drive.Stop()
	drive.Disconnect()

	// todo something with this
	lidar.Disconnect()

	audio.Sound("deactivated")
}

// Main control loop
func loop() {
	// fmt.Println("Start control loop...")
	for {
		var left = 0
		var right = 0

		var leftInRange = between(xbox.LeftY(), MinJoystickSpeed, xbox.AnalogMax) || between(xbox.LeftY(), xbox.AnalogMin, -MinJoystickSpeed)
		var rightInRange = between(xbox.RightY(), MinJoystickSpeed, xbox.AnalogMax) || between(xbox.RightY(), xbox.AnalogMin, -MinJoystickSpeed)

		if leftInRange {
			left = mapRange(xbox.LeftY(), xbox.AnalogMin, xbox.AnalogMax, drive.DriveMin, drive.DriveMax)
		}

		if rightInRange {
			right = mapRange(xbox.RightY(), xbox.AnalogMin, xbox.AnalogMax, drive.DriveMin, drive.DriveMax)
		}

		if leftInRange || rightInRange {
			drive.Drive(left, right)
		} else {
			drive.Stop()
		}

		// TODO: this should be a timer check to maintain a maximum of 60 Hz
		time.Sleep(time.Second / LoopRate)
	}
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

	// Xbox init
	xbox.Connect()
	go xbox.Control()

	lidar.Connect()
	go lidar.SocketServer()
	go lidar.Read()

	// Start drive
	drive.Connect()

	loop()

}
