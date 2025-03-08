package kafka_module

import (
	"context"
	"os"

	"encoding/json"
	"fmt"

	"net/smtp"
	"time"

	"fyp.com/m/common"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func ConsumerwithShutdown(ctx context.Context) {
	//Creating the consumer
	config := &kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT")),
		"group.id":          "my-second-app",
		"auto.offset.reset": "earliest",
	}
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		fmt.Println("error connecting to bootstrap server")
		return
	}
	defer func() {
		fmt.Println("Closing Kafka consumer...")
		consumer.Close()
	}()
	//Subscrbibe it to a topic
	err = consumer.Subscribe("update-emails", nil)
	if err != nil {
		fmt.Println("error subscribing topic")
		return
	}
	fmt.Println("Consumer started. Listening for messages...")
	//Start a channel in main thread that listens for an interupt signal and create a context for clean shutdown
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutdown signal received. Closing consumer...")
			return
		default:
			msg, err := consumer.ReadMessage(10 * time.Second)
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.Code() == kafka.ErrTimedOut {
					continue
				}
				fmt.Println("Error reading message:", err)
				continue
			}

			var emailData common.EmailMessage
			err = json.Unmarshal(msg.Value, &emailData)
			if err != nil {
				fmt.Println("Error parsing JSON:", err)
				continue
			}
			fmt.Printf("Email: %v", emailData)
			err = sendEmail(emailData)
			if err != nil {
				fmt.Println("\nError Sending message:", err)
			} else {
				fmt.Printf("\nEmail sent to %s succesfully! \n", emailData.Email)
			}

		}
	}
}

// Send email using SMTP
func sendEmail(cemailData common.EmailMessage) error {
	// SMTP Server Configuration
	smtpHost := "smtp.mailersend.net"
	smtpPort := "587"

	// Load credentials from environment variables for security
	senderEmail := "MS_JBzwfA@trial-3zxk54ve2pxljy6v.mlsender.net"
	senderPassword := "mssp.1TkYvoQ.0p7kx4xqkkvg9yjr.fK7JB7O"

	if senderEmail == "" || senderPassword == "" {
		return fmt.Errorf("SMTP credentials are missing")
	}

	// Construct email message
	message := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		senderEmail, cemailData.Email, cemailData.Subject, cemailData.Body,
	)

	// SMTP Authentication
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	// Send email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{cemailData.Email}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	fmt.Println("Email sent successfully!")
	return nil
}
