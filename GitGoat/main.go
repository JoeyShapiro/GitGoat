package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/gousb"
)

func main() {
	// Initialize a new Context.
	ctx := gousb.NewContext()
	defer ctx.Close()

	// get all devices conencted
	devices, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return true
	})
	if err != nil {
		log.Fatalf("OpenDevices: %v", err)
	}

	for _, device := range devices {
		fmt.Println(device.Product())
		fmt.Printf("%+v\n", device)

		device.Close()
	}
	return

	// Open any device with a given VID/PID using a convenience function.
	dev, err := ctx.OpenDeviceWithVIDPID(0x046d, 0xc526)
	if err != nil {
		log.Fatalf("Could not open a device: %v", err)
	}
	defer dev.Close()

	// Claim the default interface using a convenience function.
	// The default interface is always #0 alt #0 in the currently active
	// config.
	intf, done, err := dev.DefaultInterface()
	if err != nil {
		log.Fatalf("%s.DefaultInterface(): %v", dev, err)
	}
	defer done()

	// Open an OUT endpoint.
	ep, err := intf.OutEndpoint(7)
	if err != nil {
		log.Fatalf("%s.OutEndpoint(7): %v", intf, err)
	}

	// Generate some data to write.
	data := make([]byte, 5)
	for i := range data {
		data[i] = byte(i)
	}

	// Write data to the USB device.
	numBytes, err := ep.Write(data)
	if numBytes != 5 {
		log.Fatalf("%s.Write([5]): only %d bytes written, returned error is %v", ep, numBytes, err)
	}
	fmt.Println("5 bytes successfully sent to the endpoint")
}

func Actualmain() {
	router := gin.Default()
	router.GET("webhook/events", webhook)

	router.RunTLS("localhost:8888", "cert.pem", "key.pem")
}

func webhook(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
	}

	event := c.Request.Header.Get("X-GitHub-Event")
	go handleHook(event, data)

	c.JSON(http.StatusAccepted, struct {
		Uuid string `json:"uuid"`
	}{Uuid: "1234"})
}

func handleHook(event string, data []byte) {
	switch event {
	case "push":
		fmt.Println("Push event")
	case "issues":
		handleIssues(data)
	case "pull_request":
		fmt.Println("Pull request event")
	}
}

func handleIssues(data []byte) error {
	var issue HookIssue
	err := json.Unmarshal(data, &issue)
	if err != nil {
		return err
	}

	switch issue.Action {
	case "opened":
		fmt.Println("Issue opened")
	case "reponed", "closed":
		fmt.Printf("\033[1;34m%s\033[0;0m \033[0;31m%s\033[0;0m \"%s\" (\033[1;34m#%d\033[0;m)\n", issue.Sender.Login, issue.Action, issue.Issue.Title, issue.Issue.Number)
	case "labeled":
		r, g, b := colorHexToRGB(issue.Label.Color)
		fmt.Printf("\033[1;34m%s\033[0;0m \033[0;31m%s\033[0;0m \"%s\" (\033[1;34m#%d\033[0;m) as \033[38;2;%d;%d;%dm%s\033[0;m\n",
			issue.Sender.Login, issue.Action, issue.Issue.Title, issue.Issue.Number, r, g, b, issue.Label.Name)
	default:
		fmt.Printf("Issue Action \"%s\" not implemented\n", issue.Action)
	}

	return nil
}

func colorHexToRGB(hex string) (red int, green int, blue int) {
	hex = strings.Replace(hex, "#", "", 1)
	i := new(big.Int)

	i.SetString(hex[0:2], 16)
	red = int(i.Int64())
	i.SetString(hex[2:4], 16)
	green = int(i.Int64())
	i.SetString(hex[4:6], 16)
	blue = int(i.Int64())

	return
}

type HookIssue struct {
	Action string `json:"action"`
	Issue  Issue  `json:"issue"`
	Sender Sender `json:"sender"`
	Label  Label  `json:"label"`
}

type Issue struct {
	Url    string `json:"url"`
	Title  string `json:"title"`
	Number int    `json:"number"`
}

type Sender struct {
	Login string `json:"login"`
}

type Label struct {
	Color       string `json:"color"`
	Default     bool   `json:"default"`
	Description string `json:"description"`
	Id          int    `json:"id"`
	Name        string `json:"name"`
	NodeId      string `json:"node_id"`
	Url         string `json:"url"`
}
