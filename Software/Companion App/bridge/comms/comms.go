package comms

import (
	"pscreenapp/config"
	"pscreenapp/utils"

	"go.bug.st/serial"
)

var BoardBlocked chan bool = make(chan bool, 10)
var FirstFrame = true

func EstablishComms() serial.Port {
	mode := &serial.Mode{
		BaudRate: config.SerialPortBaudRate,
	}
	port, err := serial.Open(config.SerialPortToUse, mode)
	utils.CheckError(err)
	return port
}

func WaitForBoardUnblockSignal(port serial.Port) {
	buff := make([]byte, 1)
	_, err := port.Read(buff)
	utils.CheckError(err)
	// fmt.Printf("Received %v bytes\n", n)
	BoardBlocked <- false
}

func SendBytes(port serial.Port, bytes []byte) {
	if FirstFrame {
		FirstFrame = false
	} else {
		<-BoardBlocked
	}
	_, err := port.Write(bytes)
	utils.CheckError(err)
	// fmt.Printf("Sent %v bytes\n", n)
	go WaitForBoardUnblockSignal(port)
}
