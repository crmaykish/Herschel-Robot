package main

import (
	"fmt"
	"math"
	"time"

	xbox "github.com/crmaykish/Xboxdrv-Golang"
	"github.com/crmaykish/herschel/drive"
)

func main() {
	fmt.Println("Starting Herschel control program...")

	go xbox.Connect()

	drive.Connect()

	for {

		if math.Abs(float64(xbox.Xbox.LeftStick.Y)) < 1000 && math.Abs(float64(xbox.Xbox.RightStick.Y)) < 1000 {
			drive.Stop()
		} else {
			drive.Drive(xbox.Xbox.LeftStick.Y/128, xbox.Xbox.RightStick.Y/128)
		}

		time.Sleep(time.Second / 60)
	}
}
