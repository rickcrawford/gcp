package pubsub

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
)

// Client subscribes and publishes to GCP PubSub
type Client struct {
	dataChan chan []byte
	client   *pubsub.Client
}

// GetProducts is a blocking call that loads next message
func (p *Client) GetProductUpdate() <-chan []byte {
	return p.dataChan
}

// Close will exit application
func (p *Client) Close() error {
	close(p.dataChan)
	return p.client.Close()
}

// NewClient creates a new pubsub client
func NewClient(projectID, topicName, subscriptionName string) (*Client, error) {
	log.Println("starting pubsub", projectID, topicName, subscriptionName)
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}

	log.Println("create topic")
	// Create a new topic with the given name.
	topic, err := client.CreateTopic(ctx, topicName)
	if err != nil {
		topic = client.Topic(topicName)
	}
	log.Println("...")

	sub, err := client.CreateSubscription(ctx, subscriptionName, pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: 10 * time.Second,
	})
	if err != nil {
		sub = client.Subscription(subscriptionName)
	}

	dataChan := make(chan []byte, 10)
	go func() {
		log.Fatal(sub.Receive(ctx, func(ctx context.Context, message *pubsub.Message) {
			buf := make([]byte, len(message.Data))
			copy(buf, message.Data)
			dataChan <- buf
			message.Ack()
		}))
	}()

	log.Println("done!")

	return &Client{
		client:   client,
		dataChan: dataChan,
	}, nil
}
