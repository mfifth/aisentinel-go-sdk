package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mfifth/aisentinel-go-sdk"
)

func main() {
	// Initialize the client
	client, err := aisentinel.NewClient(aisentinel.Config{
		APIKey: "your-api-key-here",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Evaluate content
	result, err := client.Evaluate(context.Background(), aisentinel.EvaluationRequest{
		Content:     "Hello, world!",
		ContentType: aisentinel.ContentTypeText,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Safety Score: %.2f\n", result.Score)
	fmt.Printf("Approved: %t\n", result.Approved)
}
