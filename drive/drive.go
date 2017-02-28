package drive

import (
	"fmt"
	"log"

	"time"

	"strconv"

	"github.com/crmaykish/herschel/audio"
	"github.com/tarm/serial"
)

const DriveMax = 255
const DriveMin = -255

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
		fmt.Println("Connected to Motor Board.")
		Connected = true
		audio.Sound("communicationson")
	}
}

// Disconnect from the serial port
func Disconnect() {
	fmt.Println("Disconnecting from Motor Board...")
	port.Flush()
	port.Close()
	Connected = false
	fmt.Println("Disconnected from Motor Board")
}

func sendSerial(message string) {
	if Connected {
		port.Write([]byte(message + "\n"))
		// TODO: Up the baudrate and use a call/response system instead of just sleeping and hoping it's long enough
		time.Sleep(time.Millisecond * 16)
	}
}

func Drive(left, right int) {
	sendSerial("FL:" + strconv.Itoa(left) + "!")
	sendSerial("FR:" + strconv.Itoa(right) + "!")
	sendSerial("BL:" + strconv.Itoa(left) + "!")
	sendSerial("BR:" + strconv.Itoa(right) + "!")
}

// Stop all motors
func Stop() {
	sendSerial("FL:0!")
	sendSerial("FR:0!")
	sendSerial("BL:0!")
	sendSerial("BR:0!")
}
