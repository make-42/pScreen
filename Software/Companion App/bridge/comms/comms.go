package comms

import (
	"pscreen/config"
	"pscreen/utils"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"

	"github.com/ztrue/tracerr"
)

var BoardBlocked chan bool = make(chan bool, 10)
var FirstFrame = true

var SerialPortToUse string

func EnumSerialDevices() []*enumerator.PortDetails {
	var err error
	detectedPorts, err := enumerator.GetDetailedPortsList()
	utils.CheckError(tracerr.Wrap(err))
	return detectedPorts
}

func EstablishComms() serial.Port {
	mode := &serial.Mode{
		BaudRate: config.Config.SerialPortBaudRate,
	}
	port, err := serial.Open(SerialPortToUse, mode)
	utils.CheckError(tracerr.Wrap(err))
	return port
}

func WaitForBoardUnblockSignal(port serial.Port) {
	buff := make([]byte, 1)
	_, err := port.Read(buff)
	utils.CheckError(tracerr.Wrap(err))
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
	utils.CheckError(tracerr.Wrap(err))
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
