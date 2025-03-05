package Publisher

import (
	"encoding/json"
	"fmt"
	"module/OrderService"
	lg "module/logger"
	"time"

	"github.com/IBM/sarama"
)

// Global Producer Instance
var producer sarama.SyncProducer

// Initialize Kafka Producer (called once)
func InitKafkaProducer(brokers []string) error {
	lg.Log.Info("Initializing kafka producer")
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	var err error
	producer, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		lg.Log.Error("failed to start Kafka producer: %v", err)
		return fmt.Errorf("failed to start Kafka producer: %v", err)
	}
	lg.Log.Info("kafka producer has been initialized Successfully")
	return nil
}

// Publish message to a given Kafka topic
func PublishMessage(topic string, message *OrderService.CompleteOrder) error {
	lg.Log.Info("Entered kafka producer")
	if producer == nil {
		lg.Log.Error("Kafka producer not initialized")
		return fmt.Errorf("Kafka producer not initialized")
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		lg.Log.Error("failed to serialize message: %v", err)
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
			lg.Log.Info("Message sent to topic %s (partition %d, offset %d)\n", topic, partition, offset)
			fmt.Printf("Message sent to topic %s (partition %d, offset %d)\n", topic, partition, offset)
			//producer.Close()
			CloseKafkaProducer()
			return nil
		}
		lg.Log.Info("Retry %d: Failed to send message: %v\n", retry+1, err)
		fmt.Printf("Retry %d: Failed to send message: %v\n", retry+1, err)
		time.Sleep(time.Duration(2^retry) * time.Second)
	}
	lg.Log.Error("failed to send message after retries")
	return fmt.Errorf("failed to send message after retries")
}

func CloseKafkaProducer() {
	lg.Log.Info("Entered the close kafka producer function !!")
	if producer != nil {
		producer.Close()
		lg.Log.Info("Producer is closed")
		fmt.Println("Producer is closed")
	} else {
		lg.Log.Error("Producer is not closed")
		fmt.Println("Producer is not closed")
	}
}
