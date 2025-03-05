package NotificationService

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"module/OrderService"

	"github.com/redis/go-redis/v9"
)

var ProdMsg []OrderService.FinalOrder

//var ids string

type NotificationResponse struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

func GetMostAccessedProducts(ps []OrderService.FinalOrder) (string, error) {
	//fmt.Println("Notification is being called")

	for _, k := range ps {
		var cnt = 0
		for _, l := range ps {
			if k.OrderDts.Product_Id == l.OrderDts.Product_Id {
				cnt++
			}
		}
		if cnt > 2 {
			ProdMsg = append(ProdMsg, k)
		}
	}

	var (
		ctx       = context.Background()
		user_id   = ""
		redisAddr = "localhost:6379"
		message   = ""
	)

	// Create Redis client
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Store a notification

	//notification := "New message from Alice!"
	for _, k := range ProdMsg {
		user_id = k.OrderDts.Product_Id
		mess, err1 := json.Marshal(k)
		if err1 != nil {
			log.Fatalf("Error storing notification: %v", err1)
			return "", err1
		}
		message = string(mess)
		err := client.LPush(ctx, user_id, message).Err()
		if err != nil {
			log.Fatalf("Error storing notification: %v", err)
			return "", err
		}
	}

	// Keep only the latest 100 notifications
	client.LTrim(ctx, user_id, 0, 99)

	// Create response
	response := NotificationResponse{
		UserID:  user_id,
		Message: message,
	}

	// Convert response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Error creating JSON response: %v", err)
		return "", err
	}

	fmt.Println("Notification stored successfully!")
	return string(jsonResponse), nil

}
