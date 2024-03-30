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
	"go.bug.st/serial"
)

func main() {
	pushBlame()
}

func Actualmain() {
	router := gin.Default()
	router.GET("webhook/events", webhook)

	router.RunTLS("localhost:8888", "cert.pem", "key.pem")
}

// github webhook stuff
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

func handleIssues(data []byte) (err error) {
	var issue HookIssue
	err = json.Unmarshal(data, &issue)
	if err != nil {
		return
	}

	switch issue.Action {
	case "opened", "reponed", "closed":
		fmt.Printf("\033[1;34m%s\033[0;0m \033[0;31m%s\033[0;0m \"%s\" (\033[1;34m#%d\033[0;m)\n", issue.Sender.Login, issue.Action, issue.Issue.Title, issue.Issue.Number)

		if issue.Action == "opened" {
			err = pushBlame()
		}
	case "labeled":
		r, g, b := colorHexToRGB(issue.Label.Color)
		fmt.Printf("\033[1;34m%s\033[0;0m \033[0;31m%s\033[0;0m \"%s\" (\033[1;34m#%d\033[0;m) as \033[38;2;%d;%d;%dm%s\033[0;m\n",
			issue.Sender.Login, issue.Action, issue.Issue.Title, issue.Issue.Number, r, g, b, issue.Label.Name)
	default:
		fmt.Printf("Issue Action \"%s\" not implemented\n", issue.Action)
	}

	return
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

// blamethrower stuff
func pushBlame() (err error) {
	var n int
	mode := &serial.Mode{
		BaudRate: 9600,
	}
	port, err := serial.Open("/dev/tty.usbmodemgitgoat1", mode)
	if err != nil {
		log.Fatal(err)
	}

	buff := make([]byte, 100)
	_, err = port.Write([]byte("B"))
	if err != nil {
		return
	}

	var data string
	for {
		n, err = port.Read(buff)
		if err != nil {
			return
		}
		data = string(buff[:n])
		fmt.Printf("%s", data)
		if data == "G" {
			break
		}
	}

	return nil
}
