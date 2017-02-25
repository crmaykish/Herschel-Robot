package main

import (
	"fmt"
	"time"

	"github.com/crmaykish/herschel/drive"
)

func main() {
	fmt.Println("Hello from Herschel")

	drive.Connect()

	for {
		drive.Go()
		time.Sleep(time.Second)
		drive.Stop()
		time.Sleep(time.Second)
	}
}
