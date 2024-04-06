package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	flog      *os.File
	serialDev string
)

func main() {
	var err error

	optStdout := flag.Bool("stdout", false, "log to stdout as well as log file")
	optCert := flag.String("cert", "cert.pem", "path to cert file")
	optKey := flag.String("key", "key.pem", "path to key file")
	optSerial := flag.String("serial", "/dev/tty.usbmodemgitgoat1", "path to serial device")
	optPort := flag.String("port", "8888", "port for the webhook")
	optNoHttps := flag.Bool("no-https", false, "do not use https")
	flag.Parse()

	// log to file instead of stdout
	flog, err = os.OpenFile("gitgoat.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		fmt.Println("error opening file:", err)
		os.Exit(1)
	}
	defer flog.Close()

	// if they want to log to stdout as well
	if *optStdout {
		gin.DefaultWriter = io.MultiWriter(flog, os.Stdout)
	} else {
		gin.DefaultWriter = flog
	}

	// get the serial device
	serialDev = *optSerial

	// actually start the server
	router := gin.Default()
	router.GET("webhook/events", webhook)

	if *optNoHttps {
		fmt.Println("Running http server on port", *optPort)
		router.Run("localhost:" + *optPort)
	} else {
		fmt.Println("Running https server on port", *optPort)
		router.RunTLS("localhost:"+*optPort, *optCert, *optKey)
	}
}
