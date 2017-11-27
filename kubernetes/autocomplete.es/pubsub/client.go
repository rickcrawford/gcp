package pubsub

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
	"github.com/rickcrawford/beacon/models"
)

// ProductUpdate message identifies items to index
type ProductUpdate struct {
	Products []models.Product `json:"products"`
}

// Client subscribes and publishes to GCP PubSub
type Client struct {
	dataChan chan []byte
	client   *pubsub.Client
}

// GetProducts is a blocking call that loads next message
func (p *Client) GetProducts() ([]models.Product, error) {
	data, ok := <-p.dataChan
	if ok {
		update := new(ProductUpdate)
		err := json.Unmarshal(data, update)
		if err != nil {
			return nil, err
		}
		return update.Products, nil
	}
	return nil, errors.New("closed")
}

// Close will exit application
func (p *Client) Close() error {
	close(p.dataChan)
	return p.client.Close()
}

// NewClient creates a new pubsub client
func NewClient(projectID, topicName, subscriptionName string) (*Client, error) {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}

	// Create a new topic with the given name.
	topic, err := client.CreateTopic(ctx, topicName)
	if err != nil {
		topic = client.Topic(topicName)
	}

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

	return &Client{
		client:   client,
		dataChan: dataChan,
	}, nil
}
