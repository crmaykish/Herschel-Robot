package main

import (
	"fmt"

	"time"

	"math"

	xbox "github.com/crmaykish/Xboxdrv-Golang"
	"github.com/crmaykish/herschel/drive"
)

func main() {
	fmt.Println("Hello from Herschel")

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
