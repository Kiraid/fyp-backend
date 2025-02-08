package kafka_module

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"time"

	"fyp.com/m/common"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jordan-wright/email"
)

func ConsumerwithShutdown(ctx context.Context) {
	//Creating the consumer
	config := &kafka.ConfigMap{
		"bootstrap.servers": "broker:9092",
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
				fmt.Println("Error Sending message:", err)
			} else {
				fmt.Printf("\nEmail sent to %s succesfully! \n", emailData.Email)
			}

		}
	}
}

// Send email using SMTP
func sendEmail(emailData common.EmailMessage) error {
	// SMTP server configuration
	smtpServer := "smtp.mailersend.net" // Use your SMTP server
	smtpPort := "587"
	senderEmail := "MS_JBzwfA@trial-3zxk54ve2pxljy6v.mlsender.net"
	senderPassword := "mssp.1TkYvoQ.0p7kx4xqkkvg9yjr.fK7JB7O"

	e := email.NewEmail()
	e.From = senderEmail
	e.To = []string{emailData.Email}
	e.Subject = emailData.Subject
	e.Text = []byte(emailData.Body)

	addr := fmt.Sprintf("%s:%s", smtpServer, smtpPort)

	// Connect to the SMTP server (plaintext connection first)
	client, err := smtp.Dial(addr)
	if err != nil {
		log.Println("SMTP Connection Error:", err)
		return err
	}
	defer client.Close()

	// Start TLS encryption
	tlsConfig := &tls.Config{
		ServerName: smtpServer,
	}

	if err = client.StartTLS(tlsConfig); err != nil {
		log.Println("TLS Handshake Error:", err)
		return err
	}

	// Authenticate
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpServer)
	if err = client.Auth(auth); err != nil {
		log.Println("Authentication Error:", err)
		return err
	}

	// Send email
	err = e.Send(addr, auth)
	if err != nil {
		log.Println("SMTP Error:", err)
		return err
	}

	fmt.Println("Email sent successfully!")
	return nil
}
