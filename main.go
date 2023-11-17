package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"real_time_monitoring_system/handler"
	"real_time_monitoring_system/producer"
	"real_time_monitoring_system/utils"
	"time"
)

func main() {
	fmt.Println("Real time monitoring system")

	// Create a new Gin router
	r := gin.Default()
	r.POST("/send-data", handler.SendMessageHandler())

	go PerformanceMetricsMonitor()

	// Run the application on port 8083
	err := r.Run(":8083")
	if err != nil {
		log.Print(err)
	}
}

func PerformanceMetricsMonitor() {
	type PerformanceMatricesDataResponse struct {
		CPUUsage    int       `json:"cpu_usage"`
		MemoryUsage int       `json:"memory_usage"`
		CurrentTime time.Time `json:"current_time"`
	}
	for {
		performanceMetrics := PerformanceMatricesDataResponse{
			CPUUsage:    rand.Intn(101),
			MemoryUsage: rand.Intn(1001),
			CurrentTime: time.Now(),
		}

		err := producer.ProduceMessage(performanceMetrics, utils.PerformanceMetrics)
		if err != nil {
			fmt.Printf("Failed to produce message: %s\n", err)
		}

		time.Sleep(2 * time.Second)
	}
}
