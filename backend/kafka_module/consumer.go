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

//Kafka Consumer For Emails, Runs in the main.go as a seperate go-routine from the server

func ConsumerwithShutdown(ctx context.Context) {
	//Creating the consumer
	// config := &kafka.ConfigMap{
	// 	"bootstrap.servers": fmt.Sprintf("%s:%s", os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT")),
	// 	"group.id":          "my-second-app",
	// 	"auto.offset.reset": "earliest",
	// }
	config := &kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT")),
		"group.id":          "my-second-app",
		"auto.offset.reset": "earliest",
		"security.protocol": os.Getenv("KAFKA_SECURITY_PROTOCOL"),
		"sasl.mechanism":    os.Getenv("KAFKA_SASL_MECHANISM"),
		"sasl.username":     os.Getenv("KAFKA_SASL_USERNAME"),
		"sasl.password":     os.Getenv("KAFKA_SASL_PASSWORD"),
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
	//Runs in a go routine, listening to a channel for shutdown signal in a infinite loop; 
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutdown signal received. Closing consumer...")
			return
		default:
			//polls every 10 second for a new message
			msg, err := consumer.ReadMessage(10 * time.Second)
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.Code() == kafka.ErrTimedOut {
					continue
				}
				fmt.Println("Error reading message:", err)
				continue
			}
			//reads the value from message and converts from byte and binds it to email Struct
			var emailData common.EmailMessage
			err = json.Unmarshal(msg.Value, &emailData)
			if err != nil {
				fmt.Println("Error parsing JSON:", err)
				continue
			}
			fmt.Printf("Email: %v", emailData)
			//Calls the sendEmail function with email data
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
	senderEmail := os.Getenv("SENDER_EMAIL")
	senderPassword := os.Getenv("SENDER_PASSWORD")

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
