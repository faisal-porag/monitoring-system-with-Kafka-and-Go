package producer

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

func ProduceMessage(data interface{}, topic string) error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:19092"})
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close()

	// Produce messages to topic
	deliveryChan := make(chan kafka.Event)

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("json.Marshal.err:", err)
		return err
	}

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          jsonData,
	}, deliveryChan)

	if err != nil {
		log.Fatal(err)
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		return m.TopicPartition.Error
	} else {
		log.Printf(
			"Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic,
			m.TopicPartition.Partition,
			m.TopicPartition.Offset,
		)
	}
	return nil
}
