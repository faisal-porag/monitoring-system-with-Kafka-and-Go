package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"real_time_monitoring_system/handler"
)

func main() {
	fmt.Println("Real time monitoring system")

	// Create a new Gin router
	r := gin.Default()
	r.POST("/send-data", handler.SendMessageHandler())

	// Run the application on port 8083
	err := r.Run(":8083")
	if err != nil {
		log.Print(err)
	}
}
