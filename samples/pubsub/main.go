package main

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
)

const projectID = "typeahead-183622"
const topicName = "updates"
const subscriptionName = "test-subscription"

func main() {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Create a new topic with the given name.
	topic, err := client.CreateTopic(ctx, topicName)
	if err != nil {
		topic = client.Topic(topicName)
	}
	defer topic.Stop()

	// // Create a new subscription to the previously created topic
	// // with the given name.
	// sub, err := client.CreateSubscription(ctx, subscriptionName, pubsub.SubscriptionConfig{
	// 	Topic:       topic,
	// 	AckDeadline: 10 * time.Second,
	// })
	// if err != nil {
	// 	sub = client.Subscription(subscriptionName)
	// }

	go func() {
		for {

			res := topic.Publish(ctx, &pubsub.Message{Data: []byte("payload")})

			id, err := res.Get(ctx)
			if err != nil {
				log.Fatal("error publishing message", err)
			}
			log.Printf("Published a message with a message ID: %s\n", id)

			<-time.After(time.Second * 10)
		}

	}()

	// err = sub.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
	// 	log.Printf("Got message: %s", m.Data)
	// 	m.Ack()
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	select {}
}
