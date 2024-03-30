package main

import (
	"errors"
	"fmt"
	"time"

	"go.bug.st/serial"
)

// blamethrower stuff
func pushBlame() (err error) {
	var n int
	mode := &serial.Mode{
		BaudRate: 9600,
	}
	port, err := serial.Open("/dev/tty.usbmodemgitgoat1", mode)
	if err != nil {
		return
	}

	buff := make([]byte, 100)
	_, err = port.Write([]byte("B"))
	if err != nil {
		return
	}

	var data string
	start := time.Now()
	for {
		n, err = port.Read(buff)
		if err != nil {
			return
		}
		data = string(buff[:n])
		fmt.Printf("%s\n", data)

		if data == "G" {
			break
		}

		if time.Since(start) > 5*time.Second {
			err = errors.New("goat timeout")
			break
		}
	}

	return
}
