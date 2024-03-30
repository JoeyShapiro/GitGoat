package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("webhook/events", webhook)

	router.RunTLS("localhost:8888", "cert.pem", "key.pem")
}
