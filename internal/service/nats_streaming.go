package service

import (
	"fmt"
	"sync"

	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

type NATSService struct {
	natsConn stan.Conn
	subMu    sync.Mutex
	subMap   map[string][]chan string
}

func NewNATS(sc stan.Conn) *NATSService {
	return &NATSService{
		natsConn: sc,
		subMap:   make(map[string][]chan string),
	}
}

func (n *NATSService) SubscribeClient(client_id string) error {
	n.subMu.Lock()
	defer n.subMu.Unlock()

	subChan := make(chan string)
	if _, ok := n.subMap[client_id]; !ok {
		n.subMap[client_id] = make([]chan string, 0)
		logrus.Printf("New sub, client %s", client_id)
	}

	n.subMap[client_id] = append(n.subMap[client_id], subChan)
	return nil
}

func (n *NATSService) UnsubscribeClient(client_id string) {
	n.subMu.Lock()
	defer n.subMu.Unlock()

	if subs, ok := n.subMap[client_id]; ok {
		for _, ch := range subs {
			close(ch)
		}

		delete(n.subMap, client_id)
		logrus.Printf("Client %s unsub", client_id)
	}
}

func (n *NATSService) NotifyNewOrder(order_id int) {
	n.subMu.Lock()
	defer n.subMu.Unlock()

	subject := "new_order"
	message := fmt.Sprintf("New order id: %d", order_id)

	hasSubscribers := false
	for client_id, subs := range n.subMap {
		if len(subs) > 0 {
			hasSubscribers = true

			logrus.Printf("Notifying subscribers for client %s", client_id)

			break
		}
	}

	if !hasSubscribers {
		logrus.Println("No active subscribers, skipping notification.")
		return
	}

	if err := n.natsConn.Publish(subject, []byte(message)); err != nil {
		logrus.Printf("Failed to publish new order notification: %v", err)
		return
	}

	logrus.Printf("New order notification published: %s", message)
}
