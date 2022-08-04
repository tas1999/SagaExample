package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type OrderConsumer struct {
	conn *kafka.Reader
	rep  Repository
	prod Produser
}

func NewConsumer(conf OrderProduserConfig, rep Repository, prod Produser) (*OrderConsumer, error) {
	host := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	conn := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{host},
		GroupID:  "consumer-group-id",
		Topic:    conf.Topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	return &OrderConsumer{
		conn: conn,
		rep:  rep,
		prod: prod,
	}, nil
}
func (c *OrderConsumer) Start() error {
	for {
		fmt.Println("Start")
		order, err := c.Consume()
		if err != nil {
			fmt.Println("failed to close batch:", err)
			continue
		}
		fmt.Println("Start 2")
		order.Status = Cancelled
		fmt.Println("Start 3")
		cust, err := c.rep.GetCustomer(order.UserId)
		fmt.Println("Start 4")
		if err != nil {
			fmt.Println("failed to close batch:", err)
		}
		if cust.Credit > order.Price {
			fmt.Println("Start 5")
			order.Status = Approved
		}

		fmt.Println("Start 6")
		err = c.prod.SendAddOrderEvent(order)
		if err != nil {
			fmt.Println("failed SendAddOrderEvent", err)
		}
		fmt.Println("Start 6")
	}
}
func (c *OrderConsumer) Consume() (Order, error) {
	m, err := c.conn.ReadMessage(context.Background())
	if err != nil {
		return Order{}, err
	}
	var ord Order
	err = json.Unmarshal(m.Value, &ord)
	if err != nil {
		return Order{}, err
	}
	return ord, nil
}
