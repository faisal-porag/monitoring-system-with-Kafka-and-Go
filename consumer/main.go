package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

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
	fmt.Printf("Monitoring Data: %s\n", message)
}

func main() {
	consumeAndMonitor("your_topic_name")
}
