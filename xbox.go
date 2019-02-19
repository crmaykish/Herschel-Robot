package xbox

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

const AnalogMax = 32767
const AnalogMin = -32768

type analogStick struct {
	X int
	Y int
}

type dpad struct {
	Left  bool
	Right bool
	Up    bool
	Down  bool
}

type xbox struct {
	LeftStick    analogStick
	RightStick   analogStick
	Dpad         dpad
	LeftBumper   bool
	RightBumper  bool
	LeftTrigger  int
	RightTrigger int
	Start        bool
	Back         bool
	Guide        bool
	A            bool
	B            bool
	X            bool
	Y            bool
}

// Output from xboxdrv program to parse for input
var rawOutput io.ReadCloser

// Last known state of the controller
var state xbox

func Connect() {
	fmt.Println("Connecting to Xbox Controller...")

	cmd := exec.Command("xboxdrv")
	rawOutput, _ = cmd.StdoutPipe()

	// TODO: error checking

	cmd.Start()

	fmt.Println("Connected to Xbox Controller.")
}

func Control() {
	fmt.Println("Reading Xbox controller input...")
	scanner := bufio.NewScanner(rawOutput)

	for scanner.Scan() {
		var line = scanner.Text()
		// Length 139 indicates data line, ignore anything else
		if len(line) == 139 {
			// Parse controller data
			state.LeftStick.X, _ = strconv.Atoi(strings.Trim(line[3:9], " "))
			state.LeftStick.Y, _ = strconv.Atoi(strings.Trim(line[13:19], " "))
			state.RightStick.X, _ = strconv.Atoi(strings.Trim(line[24:30], " "))
			state.RightStick.Y, _ = strconv.Atoi(strings.Trim(line[34:40], " "))
		} else if strings.Contains(line, "[ERROR]") {
			fmt.Println(line)
		}
	}

	fmt.Println("Stopped reading Xbox controller input.")
}

func LeftY() int {
	return state.LeftStick.Y
}

func RightY() int {
	return state.RightStick.Y
}
