package NotificationService

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
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

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("✅ Email sent successfully!")
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

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("✅ Email sent successfully!")
}
