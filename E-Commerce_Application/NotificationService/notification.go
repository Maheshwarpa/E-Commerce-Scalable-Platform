package NotificationService

import (
	"fmt"
	"log"
	"module/OrderService"
	"module/UserService"
	lg "module/logger"
	"net/smtp"
	"os"
	"strconv"
	"time"
)

func SendSuccessEmail() {

	from := "challengingperson97@gmail.com"
	password := os.Getenv("EMAIL_PASSWORD")
	to := []string{"challengingperon97@example.com"}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	orderNumber := "123456"
	customerName := "John Doe"
	orderDate := "2025-02-19"

	message := []byte(fmt.Sprintf("Order Confirmation–%s\n\nDear %s,\nThank you for your order! We are pleased to confirm that your order #%s has been successfully placed on %s.\n\nThanks & Regards,\nMaheshwar PA.",
		orderNumber, customerName, orderNumber, orderDate))
	auth := smtp.PlainAuth("", from, password, smtpHost)

	err1 := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err1 != nil {
		log.Fatal(err1)
	}

	log.Println("✅ Success Email sent successfully!")
}

func SendFailEmail() {

	from := "challengingperson97@gmail.com"
	password := os.Getenv("EMAIL_PASSWORD")
	to := []string{"challengingperon97@example.com"}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	orderNumber := "123456"
	customerName := "John Doe"
	orderDate := "2025-02-19"
	reason := "Payment Failure"
	refundTime := "3-5 business days"

	message := []byte(fmt.Sprintf("Subject: Order Placement Failed – %s\n\nDear %s,\n\nWe regret to inform you that your order #%s placed on %s could not be processed successfully.\n\nReason for Failure: %s\n\nYou may try placing the order again or contact our support team for assistance. If any amount was deducted, it will be refunded to your original payment method within %s.\n\nFor further assistance, please reach out to us at support@example.com or call us at +1-800-123-4567.\n\nWe apologize for the inconvenience and appreciate your patience.\n\nBest Regards,\nMaheshwar PA.",
		orderNumber, customerName, orderNumber, orderDate, reason, refundTime))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err1 := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err1 != nil {
		log.Fatal(err1)
	}

	log.Println("✅ Fail Email sent successfully!")
}

func SendOrderConfirmationEmail(od OrderService.OrderResponse, ud *UserService.UserDetails) error {
	lg.Log.Info("Sending order confirmation email")
	fmt.Println(ud.Cust_Email)
	from := "challengingperson97@gmail.com"
	password := os.Getenv("EMAIL_PASSWORD")
	to := []string{ud.Cust_Email}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	orderNumber := od.OrderId
	customerName := ud.Cust_Name
	orderDate := od.OrderDate
	k, _ := OrderService.CalculateTotal((*OrderService.CompleteOrder)(&od))
	totalAmount := "$" + strconv.FormatFloat(k, 'f', 2, 64)

	message := []byte(fmt.Sprintf(
		"Subject: Order Confirmation - #%s\n\n"+
			"Dear %s,\n\n"+
			"Thank you for your order! 🎉 We are pleased to confirm that your order #%s has been successfully placed on %s.\n\n"+
			"Order Summary:\n"+
			"Order Number: %s\n"+
			"Total Amount: %s\n\n"+
			"We will notify you once your order has been shipped. If you have any questions, feel free to reach out.\n\n"+
			"Thanks & Regards,\n"+
			"Maheshwar PA.",
		orderNumber, customerName, orderNumber, orderDate, orderNumber, totalAmount))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Fatal(err)
		return err
	}
	lg.Log.Info("✅ Order confirmation email sent successfully!")
	log.Println("✅ Order confirmation email sent successfully!")
	return nil
}

func SendInsufficientMail(od OrderService.OrderResponse, ud *UserService.UserDetails) error {
	lg.Log.Info("Sending insufficient Email")
	fmt.Println(ud.Cust_Email)
	from := "challengingperson97@gmail.com"
	password := os.Getenv("EMAIL_PASSWORD")
	to := []string{ud.Cust_Email}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	orderNumber := od.OrderId
	customerName := ud.Cust_Name
	orderDate := od.OrderDate
	reason := "Insufficient Balance"
	refundTime := "3-5 business days"

	message := []byte(fmt.Sprintf("Subject: Order Placement Failed – %s\n\nDear %s,\n\nWe regret to inform you that your order #%s placed on %s could not be processed successfully.\n\nReason for Failure: %s\n\nYou may try placing the order again or contact our support team for assistance. If any amount was deducted, it will be refunded to your original payment method within %s.\n\nFor further assistance, please reach out to us at support@skht.com or call us at +1-800-123-4567.\n\nWe apologize for the inconvenience and appreciate your patience.\n\nBest Regards,\nMaheshwar PA.",
		orderNumber, customerName, orderNumber, orderDate, reason, refundTime))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Fatal(err)
		return err
	}
	lg.Log.Info("✅ Insufficient Email sent successfully!")
	log.Println("✅ Insufficient Email sent successfully!")
	return nil
}

func SendShippedEmail(od OrderService.OrderResponse, ud *UserService.UserDetails) error {
	lg.Log.Info("Sending shipment email !!!")
	fmt.Println(ud.Cust_Email)
	from := "challengingperson97@gmail.com"
	password := os.Getenv("EMAIL_PASSWORD")
	to := []string{ud.Cust_Email}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	orderNumber := od.OrderId
	customerName := ud.Cust_Name
	trackingNumber := "TRACK" + od.OrderId
	orderdate := time.Now()
	estimatedDelivery := orderdate.AddDate(0, 0, 3)

	message := []byte(fmt.Sprintf(
		"Subject: Your Order #%s Has Shipped! 🚚\n\n"+
			"Dear %s,\n\n"+
			"Great news! Your order #%s has been shipped. You can track your package using the tracking number: %s.\n\n"+
			"Estimated delivery date: %s.\n\n"+
			"Thank you for shopping with us!\n\n"+
			"Best Regards,\n"+
			"Maheshwar PA.",
		orderNumber, customerName, orderNumber, trackingNumber, estimatedDelivery))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Fatal(err)
		return err
	}
	lg.Log.Info("✅ Shipping email sent successfully!")
	log.Println("✅ Shipping email sent successfully!")
	return nil
}
