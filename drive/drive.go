package drive

import (
	"fmt"
	"log"

	"time"

	"github.com/tarm/serial"
)

// Connected is the state of the serial port
var Connected = false

var port *serial.Port

// Connect to the serial port
func Connect() {
	fmt.Println("Connecting to Motor Board...")
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 38400}
	var err error
	port, err = serial.OpenPort(c)

	if err != nil {
		log.Fatal(err)
	} else {
		Connected = true
	}
}

// Disconnect from the serial port
func Disconnect() {
	fmt.Println("Disconnecting from Motor Board...")
	port.Flush()
	port.Close()
	Connected = false
}

func sendSerial(message string) {
	if Connected {
		fmt.Println("sending: " + message)
		port.Write([]byte(message + "\n"))
		// TODO: Up the baudrate and use a call/response system instead of just sleeping and hoping it's long enough
		time.Sleep(time.Millisecond * 10)
	}
}

// Stop all motors
func Stop() {
	sendSerial("FL:0!")
	sendSerial("FR:0!")
	sendSerial("BL:0!")
}

func Go() {
	sendSerial("FL:255!")
	sendSerial("FR:255!")
	sendSerial("BL:255!")
}
