package comms

import (
	"pscreenapp/config"
	"pscreenapp/utils"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

var BoardBlocked chan bool = make(chan bool, 10)
var FirstFrame = true

var SerialPortToUse string

func EnumSerialDevices() []*enumerator.PortDetails {
	var err error
	detectedPorts, err := enumerator.GetDetailedPortsList()
	utils.CheckError(err)
	return detectedPorts
}

func EstablishComms() serial.Port {
	mode := &serial.Mode{
		BaudRate: config.Config.SerialPortBaudRate,
	}
	port, err := serial.Open(SerialPortToUse, mode)
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

func SetDefaultSerialPort() {
	serialDevices := EnumSerialDevices()
	for _, port := range serialDevices {
		if port.SerialNumber == config.Config.DefaultPortSerialNumber {
			SerialPortToUse = port.Name
		}
	}
}
