package Publisher

import (
	"encoding/json"
	"fmt"
	"log"

	"module/OrderService"
	"sync"

	"github.com/IBM/sarama"
)

var Ipnd *int

func PublishOrderCreatedEvent(brokers []string, topic string, message *OrderService.CompleteOrder, wg *sync.WaitGroup) {
	// Configure Kafka producer
	fmt.Println("Entered OrderCReated Event")
	defer wg.Done()
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	// Create a new Kafka producer
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to start producer: %v", err)
	}
	defer producer.Close()

	// Convert struct to JSON
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Failed to serialize message: %v", err)
	}

	// Create Kafka message
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(jsonMessage),
	}

	// Send message to Kafka
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)

}

func PublishPaymentStatusEvent(brokers []string, topic string, message *OrderService.OrderResponse, wg *sync.WaitGroup) error {
	// Configure Kafka producer
	fmt.Println("Entered Payment Producer function")
	wg.Add(1)
	defer wg.Done()
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	// Create a new Kafka producer
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to start producer: %v", err)
		return err
	}
	defer producer.Close()

	// Convert struct to JSON
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Failed to serialize message: %v", err)
		return err
	}

	// Create Kafka message
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(jsonMessage),
	}

	// Send message to Kafka
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
		return err
	}

	fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)
	fmt.Println("Before", *Ipnd)

	*Ipnd -= 1
	fmt.Println("After in Publiser", *Ipnd)

	return nil
}
