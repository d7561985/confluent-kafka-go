package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <broker> <topic>\n",
			os.Args[0])
		os.Exit(1)
	}

	broker := os.Args[1]
	topic := os.Args[2]

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Producer %v\n", p)

	done_chan := make(chan bool)

	go func() {
	outer:
		for e := range p.Events {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Err != nil {
					fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Err)
				} else {
					fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
				break outer

			default:
				fmt.Printf("Ignored event: %s\n", ev)
			}
		}

		close(done_chan)
	}()

	value := "Hello Go!"
	p.ProduceChannel <- &kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.KAFKA_PARTITION_ANY}, Value: []byte(value)}

	// wait for delivery report goroutine to finish
	_ = <-done_chan

	p.Close()
}