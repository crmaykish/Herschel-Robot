package main

import (
	"fmt"

	"github.com/crmaykish/herschel/pkg/lidar"
)

func main() {
	fmt.Println("lidar")

	lidar.Connect()
	go lidar.SocketServer()
	lidar.Read()
}
