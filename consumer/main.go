package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"os"
	"os/signal"
	"real_time_monitoring_system/utils"
	"syscall"
	"time"
)

type ConsumeMessageDataResponse struct {
	Message     string    `json:"message"`
	MessageTime time.Time `json:"message_time"`
}

type PerformanceMatricesDataResponse struct {
	CPUUsage    int       `json:"cpu_usage"`
	MemoryUsage int       `json:"memory_usage"`
	CurrentTime time.Time `json:"current_time"`
}

func consumeAndMonitor(topics []string) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:19092",
		"group.id":          "consumer_group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		log.Fatal(err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case sig := <-signalChan:
				fmt.Printf("Caught signal %v: terminating\n", sig)
				os.Exit(1)
			default:
				msg, err := consumer.ReadMessage(-1)
				if err == nil {
					processMessage(msg)
				} else {
					log.Printf("Error consuming message: %v (%v)\n", err, msg)
				}
			}
		}
	}()

	// Block until a signal is received
	select {}
}

func processMessage(msg *kafka.Message) {
	topic := *msg.TopicPartition.Topic
	switch topic {
	case utils.PerformanceMetrics:
		processPerformanceMetrics(msg.Value)
	case utils.NotificationBulk:
		processAndDisplayMonitoringData(msg.Value)
	// Add more cases for additional topics
	default:
		log.Printf("Received message from unknown topic: %s\n", topic)
	}
}

func processAndDisplayMonitoringData(message []byte) {
	var receiveMessage ConsumeMessageDataResponse
	err := json.Unmarshal(message, &receiveMessage)
	if err != nil {
		log.Println("json.Unmarshal.err:", err)
		return
	}

	fmt.Printf(
		"Notification: Message: %s | Received At: %v \n",
		receiveMessage.Message,
		utils.DateFormat(receiveMessage.MessageTime),
	)
}

func processPerformanceMetrics(message []byte) {
	var receiveMessage PerformanceMatricesDataResponse
	err := json.Unmarshal(message, &receiveMessage)
	if err != nil {
		log.Println("json.Unmarshal.err:", err)
		return
	}

	fmt.Println("")
	fmt.Printf(
		"System Monitoring Log:\nCPU Usage: %d%%\nMemory Usage: %d%%\nCurrent Time: %v\n",
		receiveMessage.CPUUsage,
		receiveMessage.MemoryUsage,
		utils.DateFormatV2(receiveMessage.CurrentTime),
	)
}

func main() {
	fmt.Println("Consumer service is running ...")
	topics := []string{
		utils.NotificationBulk,
		utils.PerformanceMetrics,
	}
	consumeAndMonitor(topics)
}
