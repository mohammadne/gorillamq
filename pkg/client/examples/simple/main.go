package main

import (
	"fmt"
	"time"

	"github.com/mohammadne/gorillamq/pkg/client"
)

func main() {
	client, err := client.NewClient("gorillamq://root:Pa$$word@localhost:8080")
	if err != nil {
		panic(err)
	}

	client.Subscribe("simple-topic", func(data []byte) {
		fmt.Println(string(data))
	})

	client.Publish("simple-topic", []byte("Simple Message"))

	time.Sleep(3 * time.Second)
}
