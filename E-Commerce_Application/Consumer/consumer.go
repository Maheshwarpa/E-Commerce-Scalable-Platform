package Consumer

/*import (
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

	"github.com/IBM/sarama"
	//"github.com/jackc/pgx/v5/pgxpool"
)

var ps OrderService.OrderResponse
var brokers = []string{"localhost:9092"}
var topics = "Payment"

var (
	pass = "SUCCESS"
	fail = "FAILED"
)*/
/*
type ConsumerHandler struct {
	wg  *sync.WaitGroup
	tpc string
}

// ConsumeMessages starts the Kafka consumer
func ConsumeMessages(brokers []string, topic string, groupID string, wgs *sync.WaitGroup) {

	fmt.Println("Entered ConsumerMessages")
	// Configure Kafka consumer
	//time.Sleep(10 * time.Second)
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// Create a consumer handler
	handler := ConsumerHandler{}
	handler.wg = wgs
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

// Setup is run at the beginning of a new session
func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	h.wg.Add(1)
	return nil
}

// Cleanup is run at the end of a session
func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	h.wg.Done()
	return nil
}

// ConsumeClaim processes messages from Kafka
func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	fmt.Println("Entered  ConsumeClaim Section")
	//time.Sleep(5 * time.Second)

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
			fmt.Println("Entered Orders")
			//defer h.wg.Done()

			flag, err1 := PaymentService.CheckEligibility(event, Database.SampleData)
			if err1 != nil {
				return fmt.Errorf("Unable to check the eligibility of payment status %s", err1)
			}
			if flag {
				//defer h.wg.Done()
				fmt.Println("At Orders Ipnd+ is", *Publisher.Ipnd)

				//Publisher.PublishPaymentStatusEvent(brokers, topics, &event, h.wg)
				//ConsumeMessages(brokers, topics, "Payment", h.wg)

				*Publisher.Ipnd -= 1
				fmt.Println("At Orders Ipnd- is", *Publisher.Ipnd)
				fmt.Println("////////Waiting for Go ROutine Orders to join")
				//h.wg.Wait()
				fmt.Println("///////////Waiting for Go ROutine Orders has joined")

			} else {
				// Send Notification to User that insufficient balance
				//defer h.wg.Done()
				fo.OrderId = event.OrderId
				fo.OrderStatus = fail
				OrderService.FinalOrderList = append(OrderService.FinalOrderList, fo)
				err2 := NotificationService.SendInsufficientMail(event, Database.SampleData)
				if err2 != nil {
					log.Fatalf("Unable to send Insufficient Balance Email %s", err2)
					return err2
				}

				*Publisher.Ipnd -= 1
				fmt.Println("At Orders Ipnd is", *Publisher.Ipnd)
			}

		} else if h.tpc == "Payment" {
			fmt.Println("Entered Payments")
			//defer h.wg.Done()
			fo.OrderId = event.OrderId
			fo.OrderStatus = pass
			OrderService.FinalOrderList = append(OrderService.FinalOrderList, fo)
			err3 := NotificationService.SendOrderConfirmationEmail(event, Database.SampleData)
			if err3 != nil {
				log.Fatalf("Unable to send Order Confirmation Email %s", err3)
				return err3
			}
			//time.Sleep(2 * time.Second)
			err4 := NotificationService.SendShippedEmail(event, Database.SampleData)
			if err4 != nil {
				log.Fatalf("Unable to send Order Shipment Confirmation Email %s", err4)
				return err4
			}
			h.wg.Done()
			//*Publisher.Ipnd -= 1
			//fmt.Println("At Payment Ipnd is", *Publisher.Ipnd)
			//h.wg.Done()
			//h.wg.Done()
			//h.wg.Wait()

		} else {
			break

		}
		//h.wg.Done()
		session.MarkMessage(msg, "")

	}
	fmt.Println("Outside of claim message")
	return nil
}
*/
