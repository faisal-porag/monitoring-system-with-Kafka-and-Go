package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func produceMessage(message string, topic string) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:19092"})
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close()

	// Produce messages to topic
	deliveryChan := make(chan kafka.Event)

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, deliveryChan)

	if err != nil {
		log.Fatal(err)
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		log.Printf(
			"Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic,
			m.TopicPartition.Partition,
			m.TopicPartition.Offset,
		)
	}
}

func sendMessage() func(c *gin.Context) {
	type Success struct {
		Status string `json:"status"`
	}

	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json; charset=utf-8")

		produceMessage("Hello, Kafka!", "your_topic_name")

		c.JSON(http.StatusOK, Success{
			Status: "Ok",
		})

		return
	}
}

func main() {
	fmt.Println("Real time monitoring system")

	// Create a new Gin router
	r := gin.Default()
	r.GET("/send-data", sendMessage())

	// Run the application on port 8083
	err := r.Run(":8083")
	if err != nil {
		log.Print(err)
	}
}
