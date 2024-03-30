package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

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
	case "created", "reponed", "closed":
		fmt.Printf("\033[1;34m%s\033[0;0m \033[0;31m%s\033[0;0m \"%s\" (\033[1;34m#%d\033[0;m)\n", issue.Sender.Login, issue.Action, issue.Issue.Title, issue.Issue.Number)

		if issue.Action == "created" {
			err = pushBlame() // <--- This is the call to the USB device
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
