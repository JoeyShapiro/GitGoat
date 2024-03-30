package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	pushBlame()
}

func Actualmain() {
	router := gin.Default()
	router.GET("webhook/events", webhook)

	router.RunTLS("localhost:8888", "cert.pem", "key.pem")
}
