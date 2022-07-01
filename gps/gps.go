package gps

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/adrianmo/go-nmea"
	"github.com/albenik/go-serial"
	"github.com/sirupsen/logrus"
)

type GPS struct {
	Port      serial.Port
	Waypoints []Waypoint
}

type Waypoint struct {
	Longitude     float64
	Latitude      float64
	UnixMicroTime int64
}

func NewGPS() GPS {
	// look for the gps dongle and open it
	mode := &serial.Mode{
		BaudRate: 9600,
		StopBits: serial.OneStopBit,
		DataBits: 8,
		Parity:   serial.NoParity, // TODO is this right?
	}
	// open a serial port to the gps dongle
	port, err := serial.Open("/dev/ttyACM0", mode)
	// if it's not found, exit with an error
	if err != nil {
		logrus.WithError(err).Fatal("Could not open serial port")
	}
	return GPS{
		Port:      port,
		Waypoints: make([]Waypoint, 0),
	}
}

func (g *GPS) GetWaypoint() (Waypoint, error) {
	buf := make([]byte, 1024)
	retries := 3
	for {
		_, err := g.Port.Read(buf)
		if err != nil {
			if strings.Contains(err.Error(), "EOF") {
				continue
			}
			return Waypoint{}, err
		}
		waypoint, err := g.Parse(string(buf))
		if err != nil {
			if retries > 0 {
				retries--
				continue
			}
			return Waypoint{}, err
		}
		return waypoint, nil
	}
}

func (g *GPS) Parse(data string) (Waypoint, error) {
	data = strings.Trim(data, "\x00")
	data = strings.TrimRight(data, "\r\n")
	sentences := strings.Split(data, "\r\n")

	for i := range sentences {
		if len(sentences[i]) == 0 || sentences[i][0] != '$' {
			continue
		}
		s, err := nmea.Parse(sentences[i])
		if err != nil {
			return Waypoint{}, err
		}
		if s.DataType() == nmea.TypeGLL {
			m := s.(nmea.GLL)
			//fmt.Println("lat:", m.Latitude, "lon:", m.Longitude)
			return Waypoint{
				Latitude:      m.Latitude,
				Longitude:     m.Longitude,
				UnixMicroTime: time.Now().UnixMicro(),
			}, nil
		}
	}
	return Waypoint{}, fmt.Errorf("could not parse data '%s'", data)
}

func (g *GPS) Close() error {
	return g.Port.Close()
}

func getProcessOwner() string {
	stdout, err := exec.Command("ps", "-o", "user=", "-p", strconv.Itoa(os.Getpid())).Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return strings.Trim(string(stdout), "\n")
}
