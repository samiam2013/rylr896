package main

import (
	"bytes"
	"strings"
	"time"

	"github.com/albenik/go-serial"
	//"github.com/samiam2013/rylr896/gps"
	"github.com/sirupsen/logrus"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
)

func main() {
	mode := &serial.Mode{
		BaudRate: 115_200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open("/dev/ttyS0", mode)
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
		"AT+PARAMETER=10,2,1,7",
		"AT+BAND=432500000", // 902300000,
		"AT+ADDRESS=1",
		"AT+NETWORKID=6",
		"AT+CRFOP=15",
	}
	for _, cmd := range cmds {
		resp := sendCommand(port, cmd)
		logrus.Infof("'%s' ran, result: '%s'", cmd, resp)
	}

	//gps := gps.NewGPS()

	// set pin 12 low and flash the LED when there's a transmission
	if _, err := host.Init(); err != nil {
		logrus.Fatalf("Failed to host.Init() for periphio: %s", err.Error())
	}

	whiteLED := rpi.P1_12
	if err := whiteLED.Out(gpio.Low); err != nil {
		logrus.WithError(err).Fatalf("Failed to init white LED low")
	}

	for {
		if n, err := port.ReadyToRead(); n > 0 && err == nil {
			time.Sleep(time.Millisecond * 20)
			readBuf := make([]byte, 0)
			for {
				buf := make([]byte, 100)
				n, err := port.Read(buf)
				if err != nil && !strings.HasSuffix(err.Error(), "interrupted system call") {
					logrus.WithError(err).Error("Could not read port")
				}
				if n == 0 {
					break
				}
				readBuf = append(readBuf, bytes.Trim(buf, "\x00")...)
			}
			//wp, err := gps.GetWaypoint()
			/*if err != nil {
				logrus.WithError(err).Error("Could not get GPS waypoint")
			}*/
			// TODO remove the newline that's being printed here ?
			//logrus.Infof("Data from port: %s (lat: %f, lon: %f, unix_micro: %d)",
			//	string(readBuf), wp.Latitude, wp.Longitude, wp.UnixMicroTime)
			whiteLED.Out(gpio.High)
			time.Sleep(250 * time.Millisecond)
			whiteLED.Out(gpio.Low)
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

	response := make([]byte, 0)
	var n2 uint32
	for {
		lastIter := n2
		n2, err = p.ReadyToRead()
		if err != nil {
			logrus.WithError(err).Error("Could not read response to command.")
		} else if n2 > 0 && lastIter == n2 {
			logrus.Infof("bytes for reading: %d", n2)
			break
		}
		time.Sleep(20 * time.Millisecond) // TODO can this be removed
	}
	for {
		buf := make([]byte, n2)
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
