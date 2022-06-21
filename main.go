package main

import (
	"github.com/albenik/go-serial"
	"github.com/sirupsen/logrus"
)

func main() {
	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: 0,
	}
	port, err := serial.Open("/dev/serial0", mode)
	if err != nil {
		logrus.WithError(err).Error("failed to open /dev/serial0")
	}

	// TODO break this into it's own function?
	n, err := port.Write([]byte("AT\r\n"))
	if err != nil {
		logrus.WithError(err).Fatal("Failed to send message")
	}
	logrus.Infof("Sent %v bytes\n", n)

	response := make([]byte, 1000)
	for {
		buf := make([]byte, 100)
		n, err := port.Read(buf)
		if err != nil {
			logrus.WithError(err).Error("Could not read port")
			break
		}
		if n == 0 {
			logrus.Info("EOF")
			break
		}
		response = append(response, buf...)
	}
	logrus.Infof("response:'%v'", string(response))

}
