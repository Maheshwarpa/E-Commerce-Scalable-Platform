package Consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"module/Database"
	"module/NotificationService"
	"module/OrderService"
	"module/PaymentService"
	"module/Publisher"
	"sync"
	"time"

	"github.com/IBM/sarama"
	//"github.com/jackc/pgx/v5/pgxpool"
)

var ps OrderService.OrderResponse
var brokers = []string{"localhost:9092"}
var topics = "Payment"

var (
	pass = "SUCCESS"
	fail = "FAILED"
)

// ConsumeMessages starts the Kafka consumer
func ConsumeMessages(brokers []string, topic string, groupID string, wg *sync.WaitGroup) {

	//fmt.Println("Entered ConsumerMessages")
	// Configure Kafka consumer
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// Create a consumer handler
	handler := ConsumerHandler{}
	handler.wg = wg
	handler.tpc = topic

	// Create consumer group
	for {
		consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
		if err != nil {
			log.Fatalf("Error creating consumer group: %v", err)
		}
		defer consumerGroup.Close()

		err1 := consumerGroup.Consume(context.Background(), []string{topic}, &handler)
		if err1 != nil {
			log.Fatalf("Error consuming messages: %v", err1)
		}

	}

}

// ConsumerHandler handles consumed messages
type ConsumerHandler struct {
	wg  *sync.WaitGroup
	tpc string
}

// Setup is run at the beginning of a new session
func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

// Cleanup is run at the end of a session
func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim processes messages from Kafka
func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	//fmt.Println("Entered  ConsumeClaim Section")
	for msg := range claim.Messages() {
		//fmt.Println("Entered claim message")
		var event OrderService.OrderResponse
		var fo OrderService.FinalOrder
		err := json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Printf("Error decoding message: %v", err)
			continue
		}
		if h.tpc == "Orders" {
			flag, err1 := PaymentService.CheckEligibility(event, Database.SampleData)
			if err1 != nil {
				return fmt.Errorf("Unable to check the eligibility of payment status %s", err1)
			}
			if flag {
				Publisher.PublishPaymentStatusEvent(brokers, topics, &event, (h.wg))
			} else {
				// Send Notification to User that insufficient balance
				fo.OrderId = event.OrderId
				fo.OrderStatus = fail
				OrderService.FinalOrderList = append(OrderService.FinalOrderList, fo)
				err2 := NotificationService.SendInsufficientMail(event, Database.SampleData)
				if err2 != nil {
					log.Fatalf("Unable to send Insufficient Balance Email %s", err2)
					return err2
				}
			}

			h.wg.Done()
		} else if h.tpc == "Payment" {
			fo.OrderId = event.OrderId
			fo.OrderStatus = pass
			OrderService.FinalOrderList = append(OrderService.FinalOrderList, fo)
			err3 := NotificationService.SendOrderConfirmationEmail(event, Database.SampleData)
			if err3 != nil {
				log.Fatalf("Unable to send Order Confirmation Email %s", err3)
				return err3
			}
			time.Sleep(2 * time.Second)
			err4 := NotificationService.SendShippedEmail(event, Database.SampleData)
			if err4 != nil {
				log.Fatalf("Unable to send Order Shipment Confirmation Email %s", err4)
				return err4
			}
			h.wg.Done()

		} else {
			h.wg.Done()
			continue

		}

		session.MarkMessage(msg, "")

	}
	//fmt.Println("Outside of claim message")
	return nil
}
