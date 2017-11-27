package pubsub

import (
	"testing"
)

func TestClient(t *testing.T) {
	// client, err := NewClient("typeahead-183622", "test", "test-subscription")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// defer client.Close()

	// go func() {
	// 	select {
	// 	case m := <-client.NextMessage():
	// 		log.Println(m.Data)
	// 		m.Ack()
	// 	case <-time.After(time.Second * 10):
	// 		return
	// 	}
	// }()

	// id, err := client.Publish(&pubsub.Message{Data: []byte("test")})
	// log.Println(id, err)
	// <-time.After(time.Second * 20)
}
