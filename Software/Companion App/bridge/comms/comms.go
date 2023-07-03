package comms

import (
	"fmt"
	"pscreenapp/config"
	"pscreenapp/utils"

	"go.bug.st/serial"
)

func EstablishComms() serial.Port {
	mode := &serial.Mode{
		BaudRate: config.SerialPortBaudRate,
	}
	port, err := serial.Open(config.SerialPortToUse, mode)
	utils.CheckError(err)
	return port
}

func SendBytes(port serial.Port, bytes []byte) {
	n, err := port.Write(bytes)
	utils.CheckError(err)
	fmt.Printf("Sent %v bytes\n", n)
}
