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
		StopBits: serial.OneStopBit,
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
	// TODO ? check if the response ^ to that was '+OK'

	cmds := []string{
		"AT+PARAMETER=12,7,1,7",
		"AT+BAND=915000000",
		"AT+ADDRESS=1",
		"AT+NETWORKID=6",
		"AT+CRFOP=15",
	}
	for _, cmd := range cmds {
		logrus.Infof("'%s' ran, result: '%s'", cmd, sendCommand(port, cmd))
	}

	for {
		if n, err := port.ReadyToRead(); n > 0 && err == nil {
			readBuf := make([]byte, 1000)
			for {
				buf := make([]byte, 100)
				n, err := port.Read(buf)
				if err != nil {
					logrus.WithError(err).Error("Could not read port")
				}
				if n == 0 {
					break
				}
				readBuf = append(readBuf, buf...)
			}
			logrus.Infof("Data from port: %s", string(readBuf))
		}
	}
}

// func sendMessage(p serial.Port, msg string) (string, error) {
// 	// TODO sanitize (or error? on new lines)
// 	response := sendCommand(p, "AT+SEND="+msg)
// 	if !strings.HasPrefix(response, "+OK") {
// 		return "", fmt.Errorf("not '+OK' sending message: %s", response)
// 	}
// 	return response, nil
// }

func sendCommand(p serial.Port, cmd string) string {
	n, err := p.Write([]byte(cmd + "\r\n"))
	if err != nil {
		logrus.WithError(err).Fatal("Failed to send message")
	}
	logrus.Infof("Sent %v bytes: '%s'\n", n, cmd)

	//time.Sleep(time.Millisecond * 1000)

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

// map strings in the pattern /\+ERR=([d]{1-2})/ based on the number to a string
