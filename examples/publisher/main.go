package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	v1 "github.com/UniDropOfficial/uninats-go/gen/unidrop/uninats/v1"
	"github.com/UniDropOfficial/uninats-go/uninats"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
)

var (
	client = uninats.New("nats://localhost:4333")
)

func init() {
	client.Connect()
}

func main() {
	defer client.Close()

	ctx := context.Background()

	// Initialize the event bus
	if err := client.CreateStream(ctx, jetstream.StreamConfig{
		Name:        "product_stream",
		Description: "Sample events",
		Subjects:    []string{"sample.>"},
		MaxBytes:    1024 * 1024 * 1024,
	}); err != nil {
		log.Println("🚨 [RIMNats]: Failed to initialize event bus:", err)
		os.Exit(1)
	}

	for {
		// Get the current time
		currentTime := time.Now()

		// List of sample words
		words := []string{"Apple", "Banana", "Orange", "Mango", "Grape", "Peach", "Plum", "Cherry", "Lemon", "Lime"}
		subjects := []string{"sample.created", "sample.updated"}

		// Generate a random word
		randomWord := words[currentTime.UnixNano()%int64(len(words))]
		subject := subjects[currentTime.UnixNano()%int64(len(subjects))]

		// Create a ProductCreated event
		event := &v1.Event{
			Name: subject,
			Product: &v1.ProductCreated{
				Id:         uuid.NewString(),
				Name:       randomWord,
				SupplierId: uuid.NewString(),
				CreatedAt:  currentTime.UnixMilli(),
			},
		}

		// Publish the event
		if err := client.Publish(ctx, subject, event); err != nil {
			log.Fatalf("🚨 Failed to publish event: %v", err)
		}

		fmt.Printf("🚀 Event published to subject: %s successfully! 🚀\n", subject)

		time.Sleep(time.Duration(3) * time.Second)
	}
}
