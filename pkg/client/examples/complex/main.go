package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mohammadne/gorillamq/pkg/client"
)

const gorillamqAddress = "gorillamq://localhost:8080"

func main() {
	topics := []string{"topic1", "topic2", "topic3", "topic4"}
	subscribers(topics, 1)
	time.Sleep(5000 * time.Millisecond)
	publishers(topics, 10)
	time.Sleep(3 * time.Second)
}

func subscribers(topics []string, subscribers int) {
	for index := 0; index < subscribers; index++ {
		go func(subscriberIndex int) {
			client, er := client.NewClient(gorillamqAddress)
			if er != nil {
				log.Fatal(er)
			}

			for _, topic := range topics {
				go client.Subscribe(topic, func(bytes []byte) {
					log.Printf("subscriber '%d', recieved message '%s' on topic: '%s'\n", subscriberIndex, string(bytes), topic)
				})
			}

			select {}
		}(index)
	}
}

func publishers(topics []string, publishers int) {
	for index := 0; index < publishers; index++ {
		go func(publisherIndex int) {
			client, er := client.NewClient(gorillamqAddress)
			if er != nil {
				log.Fatal(er)
			}

			for _, topic := range topics {
				msg := fmt.Sprintf("publisher %d", publisherIndex)
				client.Publish(topic, []byte(msg))
			}

			select {}
		}(index)
	}
}
