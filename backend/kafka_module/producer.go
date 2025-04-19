package kafka_module

import (
	"encoding/json"
	"fmt"

	"fyp.com/m/common"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)
//Producer for producing a message in the topic
func Producer(value common.EmailMessage) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": "kafka-broker:9092",
	}
	//Initialize the producer
	producer, err := kafka.NewProducer(config)
	if err != nil {
		panic(fmt.Sprintf("Failed to create producer: %s", err))
	}

	defer producer.Close()
	// Validity check for message delivarance 
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\nTimeStamp: %v\nOffset: %v\n", ev.TopicPartition, ev.Timestamp, ev.TopicPartition.Offset)
				}
			}
		}
	}()
	//Converts email data to bytes
	jsonValue, err := json.Marshal(value)
	if err != nil {
		fmt.Printf("Error marshalling JSON: %s\n", err)
		return
	}
	//Meta data for the producer to write
	topic := "update-emails"
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          jsonValue,
	}
	//Produce the message
	err = producer.Produce(message, nil)

	if err != nil {
		fmt.Printf("Failed to produce message: %s\n", err)
	} else {
		fmt.Println("Message sent successfully!")
	}

	// Wait for message deliveries to complete
	producer.Flush(15 * 1000)

}
