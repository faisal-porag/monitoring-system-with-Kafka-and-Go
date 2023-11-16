package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"real_time_monitoring_system/utils"
	"time"
)

type ConsumeMessageDataResponse struct {
	Message     string    `json:"message"`
	MessageTime time.Time `json:"message_time"`
}

func consumeAndMonitor(topic string) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:19092",
		"group.id":          "your_consumer_group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatal(err)
	}

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			processAndDisplayMonitoringData(msg.Value)
		} else {
			log.Printf("Error consuming message: %v (%v)\n", err, msg)
		}
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

func main() {
	consumeAndMonitor("your_topic_name")
}
