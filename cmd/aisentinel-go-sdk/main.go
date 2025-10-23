package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mfifth/aisentinel-go-sdk/pkg/aisentinel"
)

func main() {
	// Initialize the client
	client, err := aisentinel.NewClient("your-api-key-here")
	if err != nil {
		log.Fatal(err)
	}

	// Evaluate content
	result, err := client.Evaluate(context.Background(), aisentinel.DecisionRequest{
		Policy: "default",
		Input:  "Hello, world!",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Decision: %s\n", result.Decision)
}
