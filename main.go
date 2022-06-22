package main

import (
	"strings"

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
	defer func() {
		err := port.Close()
		if err != nil {
			logrus.WithError(err).Error("Could not close serial port")
		}
	}()

	logrus.Infof("response:'%v'", sendCommand(port, "AT"))
	logrus.Infof("response:'%v'", sendCommand(port, "AT+IPR?"))
}

func sendCommand(p serial.Port, cmd string) string {
	n, err := p.Write([]byte(cmd + "\r\n"))
	if err != nil {
		logrus.WithError(err).Fatal("Failed to send message")
	}
	logrus.Infof("Sent %v bytes\n", n)

	response := make([]byte, 1000)
	for {
		buf := make([]byte, 100)
		n, err := p.Read(buf)
		if err != nil {
			if strings.Contains(err.Error(), "interrupted") {
				continue
			}
			logrus.WithError(err).Error("Could not read port")
			break
		}
		if n == 0 {
			logrus.Info("EOM")
			break
		}
		response = append(response, buf...)
	}
	return string(response)
}
