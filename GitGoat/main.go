package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// log to file instead of stdout
	f, _ := os.OpenFile("gin.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	gin.DefaultWriter = f

	router := gin.Default()
	router.GET("webhook/events", webhook)

	router.RunTLS("localhost:8888", "cert.pem", "key.pem")
}
