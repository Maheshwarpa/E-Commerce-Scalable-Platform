package Consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"module/Cassandra"
	"module/Database"
	"module/NotificationService"
	"module/OrderService"
	"module/PaymentService"
	"module/Publisher"
	lg "module/logger"
	"sync"

	"github.com/IBM/sarama"
)

const (
	passs = "SUCCESS"
	faill = "FAILURE"
)

type ConsumerHandler struct {
	wg  *sync.WaitGroup
	tpc string
}

var brokers = []string{"localhost:9092"}

// ConsumeMessages starts the Kafka consumer
func ConsumeMessages(brokers []string, topic string, groupID string, wgs *sync.WaitGroup) {
	lg.Log.Info("Entered ConsumerMessages", topic)
	fmt.Println("Entered ConsumerMessages", topic)
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

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		lg.Log.Error("Error creating consumer group: %v", err)
		log.Fatalf("Error creating consumer group: %v", err)
	}
	//defer consumerGroup.Close()
	err1 := consumerGroup.Consume(context.Background(), []string{topic}, &handler)
	if err1 != nil {
		lg.Log.Error("Error consuming messages: %v", err1)
		log.Fatalf("Error consuming messages: %v", err1)
	}
	consumerGroup.Close()

}

// ConsumerHandler handles consumed messages

// Setup is run at the beginning of a new session
func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {

	return nil
}

// Cleanup is run at the end of a session
func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {

	return nil
}

// ConsumeClaim processes messages from Kafka
func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	lg.Log.Info("Entered  ConsumeClaim Section", h.tpc)
	fmt.Println("Entered  ConsumeClaim Section", h.tpc)
	//time.Sleep(5 * time.Second)
	//defer h.wg.Done()
	for msg := range claim.Messages() {
		//fmt.Println("Entered claim message")
		var event OrderService.OrderResponse
		//var fo OrderService.FinalOrder
		//fmt.Println(len(msg.Value))
		err := json.Unmarshal(msg.Value, &event)
		if err != nil {
			lg.Log.Error("Error decoding message: %v", err)
			log.Printf("Error decoding message: %v", err)
			continue
		}

		fmt.Println(event)
		processMessage(h.tpc, event, h.wg)
		//
		lg.Log.Info("Before Done C-Claim: WaitGroup count: %d\n", h.wg)
		//fmt.Printf("Before Done C-Claim: WaitGroup count: %d\n", h.wg)
		//h.wg.Done()
		lg.Log.Info("After Done C-Claim: WaitGroup count: %d\n", h.wg)
		//fmt.Printf("After Done C-Claim: WaitGroup count: %d\n", h.wg)
		session.MarkMessage(msg, "")
		break
	}
	lg.Log.Info("Outside of claim message", h.tpc)
	//fmt.Println("Outside of claim message", h.tpc)
	return nil
}

// Process messages based on topic
func processMessage(topic string, msg OrderService.OrderResponse, wg *sync.WaitGroup) {
	lg.Log.Info("Entered process message function")
	var order OrderService.OrderResponse
	order = msg
	var fo OrderService.FinalOrder
	if topic == "orders" {
		lg.Log.Info("Entered Orders")
		//fmt.Println("Entered Orders")
		//var order OrderService.OrderResponse
		//var fo OrderService.FinalOrder
		lg.Log.Info("Processing Order:", order.OrderId)
		//fmt.Println("Processing Order:", order.OrderId)

		flag, err1 := PaymentService.CheckEligibility(order, Database.SampleData)
		if err1 != nil {
			lg.Log.Error("Unable to check the eligibility of payment status %s", err1)
			fmt.Errorf("Unable to check the eligibility of payment status %s", err1)
			return
		}
		if flag {
			// Trigger producer to send message to "payment" topic
			paymentMessage := OrderService.CompleteOrder{
				OrderId:     order.OrderId,
				PlacedOrder: order.PlacedOrder,
				OrderDate:   order.OrderDate,
			}
			//wg.Add(1)
			go func() {
				lg.Log.Info("Before Done: WaitGroup count: %d\n", wg)
				//fmt.Printf("Before Done: WaitGroup count: %d\n", wg)
				defer wg.Done()
				lg.Log.Info("After Done: WaitGroup count: %d\n", wg)
				//fmt.Printf("After Done: WaitGroup count: %d\n", wg)
				ConsumeMessages(brokers, "payment", "PaymentGroup", wg)
			}()
			Publisher.InitKafkaProducer(brokers)
			//go ConsumeMessages(brokers, "payment", "Payment", wg)
			err := Publisher.PublishMessage("payment", &paymentMessage)
			if err != nil {
				lg.Log.Info("Failed to send payment message: %v", err)
				log.Printf("Failed to send payment message: %v", err)
				return
			}

			wg.Wait()
			//wg.Done()
			//return
		} else {
			fo.OrderId = order.OrderId
			fo.OrderDts = order.PlacedOrder
			fo.OrderDte = order.OrderDate
			fo.OrderStatus = faill
			OrderService.FinalOrderList = append(OrderService.FinalOrderList, fo)
			err12 := Cassandra.AddFinalOrder(Cassandra.Session, fo)
			if err12 != nil {
				lg.Log.Error("Unable to load final order into the final order table %s", err12)
				log.Fatalf("Unable to load final order into the final order table %s", err12)
				return
			}
			_, err23 := OrderService.UpdateStockInDb(order.PlacedOrder.Product_Id, order.PlacedOrder.Count, fo.OrderStatus)
			if err23 != nil {
				lg.Log.Error("Unable to Update the Balance in Product Table %s", err23)
				log.Fatalf("Unable to Update the Balance in Product Table %s", err23)
				return
			}
			err2 := NotificationService.SendInsufficientMail(order, Database.SampleData)
			if err2 != nil {
				lg.Log.Error("Unable to send Insufficient Balance Email %s", err2)
				log.Fatalf("Unable to send Insufficient Balance Email %s", err2)
				return
			}

			wg.Done()
			//return

		}

	} else if topic == "payment" {
		var payment OrderService.OrderResponse
		payment = msg
		lg.Log.Info("Entered Payments")
		//fmt.Println("Entered Payments")

		fo.OrderId = payment.OrderId
		fo.OrderDts = payment.PlacedOrder
		fo.OrderDte = payment.OrderDate
		fo.OrderStatus = passs
		OrderService.FinalOrderList = append(OrderService.FinalOrderList, fo)
		err12 := Cassandra.AddFinalOrder(Cassandra.Session, fo)
		if err12 != nil {
			lg.Log.Error("Unable to load final order into the final order table %s", err12)
			log.Fatalf("Unable to load final order into the final order table %s", err12)
			return
		}
		_, err23 := OrderService.UpdateStockInDb(payment.PlacedOrder.Product_Id, payment.PlacedOrder.Count, fo.OrderStatus)
		if err23 != nil {
			lg.Log.Error("Unable to Update the Balance in Product Table %s", err23)
			log.Fatalf("Unable to Update the Balance in Product Table %s", err23)
			return
		}

		err3 := NotificationService.SendOrderConfirmationEmail(payment, Database.SampleData)
		if err3 != nil {
			lg.Log.Error("Unable to send Order Confirmation Email %s", err3)
			log.Fatalf("Unable to send Order Confirmation Email %s", err3)
			return
		}
		//time.Sleep(2 * time.Second)
		err4 := NotificationService.SendShippedEmail(payment, Database.SampleData)
		if err4 != nil {
			lg.Log.Error("Unable to send Order Shipment Confirmation Email %s", err4)
			log.Fatalf("Unable to send Order Shipment Confirmation Email %s", err4)
			return
		}
		lg.Log.Info("Processing Payment for Order:", payment.OrderId)
		//fmt.Println("Processing Payment for Order:", payment.OrderId)
		//wg.Done()
	}
}
