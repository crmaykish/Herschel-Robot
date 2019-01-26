package lidar

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/tarm/serial"
)

var port *serial.Port
var reader *bufio.Reader

var serialConnected bool

var client net.Conn

var clientConnected bool

type LidarData struct {
	Distance        uint16
	Invalid         bool
	StrengthWarning bool
	SignalStrength  uint16
}

type Packet struct {
	Start    uint8
	Index    uint8
	Speed    uint16
	Data     [4]LidarData
	Checksum uint16
}

func parsePacket(a [22]byte) Packet {
	var packet Packet

	// Start packet should always be 0xFA, TODO: check that
	packet.Start = a[0]

	// Index byte goes from 0xA0 (packet 0, reading 0 - 3) to 0xF9 (packet 89, readings 356 - 359)
	packet.Index = a[1] - 0xA0

	// Rotational speed in RPM, TODO: I think there's more precision available here than just dividing by 64
	// TODO: Add a PID control loop to maintain constant speed
	packet.Speed = ((uint16(a[3]) << 8) + uint16(a[2])) / 64

	// Parse the four data readings
	for i := 0; i < 4; i++ {
		b0 := uint32(a[4*i+4])
		b1 := uint32(a[4*i+5])
		b2 := uint32(a[4*i+6])
		b3 := uint32(a[4*i+7])

		// Distance
		packet.Data[i].Distance = uint16(((b1 & 0x00111111) << 8) + b0)

		// Invalid flag
		packet.Data[i].Invalid = ((b1 & 0x10000000) >> 7) == 1

		// Strength warning flag
		packet.Data[i].StrengthWarning = ((b1 & 0x01000000) >> 6) == 1

		// Signal strength
		packet.Data[i].SignalStrength = uint16((b3 << 8) + b2)
	}

	// TODO: add the checksum calculation

	return packet
}

func printPacket(p Packet) {
	fmt.Printf("%X %d %d\n", p.Start, p.Index, p.Speed)
	for i := 0; i < 4; i++ {
		fmt.Printf("%t %t %d %d\n", p.Data[i].Invalid, p.Data[i].StrengthWarning, p.Data[i].Distance, p.Data[i].SignalStrength)
	}
}

func printPacketCSV(p Packet) string {
	var r string
	for i := 0; i < 4; i++ {
		r += fmt.Sprintf("%d,%d,%d\n", int(p.Index)*4+i, p.Data[i].Distance, p.Data[i].SignalStrength)
	}
	return r
}

func SocketServer() {
	fmt.Println("Starting socket server...")
	ln, _ := net.Listen("tcp", ":9000")

	clientConnected = false

	for {
		client, _ = ln.Accept()
		fmt.Println("Client connected")
		clientConnected = true
	}

}

func Connect() {
	fmt.Println("Connecting to Lidar...")
	c := &serial.Config{Name: "/dev/ttyUSB1", Baud: 115200}
	var err error
	port, err = serial.OpenPort(c)

	if err != nil {
		log.Fatal(err)
	} else {
		reader = bufio.NewReader(port)
		serialConnected = true
	}
}

// Disconnect from the serial port
func Disconnect() {
	fmt.Println("Disconnecting from Lidar...")
	port.Flush()
	port.Close()
	serialConnected = false
	fmt.Println("Disconnected from Lidar")
}

func Read() {
	var r [22]byte
	var i int

	fmt.Println("Reading LIDAR data...")

	for {
		if serialConnected {
			b, err := reader.ReadByte()

			if err != nil {
				log.Print(err)
			} else {

				if b == 0xFA {
					// clear byte array
					for j := 0; j < 22; j++ {
						r[j] = 0
					}
					i = 0

				}
				r[i] = b
				i++

				if i == 22 {
					// got there
					packet := parsePacket(r)

					for i := 0; i < 4; i++ {
						fmt.Printf("%d, %d, %d\n", int(packet.Index)*4+i, packet.Data[i].Distance, packet.Data[i].SignalStrength)

						if clientConnected {
							client.Write([]byte(fmt.Sprintf("*%d,%d,%d!\n", int(packet.Index)*4+i, packet.Data[i].Distance, packet.Data[i].SignalStrength)))
						}
					}

					// clear byte array
					for j := 0; j < 22; j++ {
						r[j] = 0
					}
					i = 0
				}
			}
		}
	}
}
