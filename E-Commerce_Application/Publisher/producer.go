package Publisher

import (
	"encoding/json"
	"fmt"
	"module/OrderService"
	"time"

	"github.com/IBM/sarama"
)

// Global Producer Instance
var producer sarama.SyncProducer

// Initialize Kafka Producer (called once)
func InitKafkaProducer(brokers []string) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	var err error
	producer, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return fmt.Errorf("failed to start Kafka producer: %v", err)
	}
	return nil
}

// Publish message to a given Kafka topic
func PublishMessage(topic string, message *OrderService.CompleteOrder) error {
	if producer == nil {
		return fmt.Errorf("Kafka producer not initialized")
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to serialize message: %v", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(jsonMessage),
	}

	// Retry mechanism with exponential backoff
	for retry := 0; retry < 3; retry++ {
		partition, offset, err := producer.SendMessage(msg)
		if err == nil {
			fmt.Printf("Message sent to topic %s (partition %d, offset %d)\n", topic, partition, offset)
			//producer.Close()
			CloseKafkaProducer()
			return nil
		}
		fmt.Printf("Retry %d: Failed to send message: %v\n", retry+1, err)
		time.Sleep(time.Duration(2^retry) * time.Second)
	}

	return fmt.Errorf("failed to send message after retries")
}

func CloseKafkaProducer() {
	if producer != nil {
		producer.Close()
		fmt.Println("Producer is closed")
	} else {
		fmt.Println("Producer is not closed")
	}
}
