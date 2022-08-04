package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type OrderProduser struct {
	conn *kafka.Conn
}
type OrderProduserConfig struct {
	Host  string `mapstructure:"host"`
	Port  int    `mapstructure:"port"`
	Topic string `mapstructure:"tppic"`
}

func NewOrderProduser(conf OrderProduserConfig) (*OrderProduser, error) {
	partition := 0
	host := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	fmt.Println(host)
	conn, err := kafka.DialLeader(context.Background(), "tcp", host, conf.Topic, partition)
	if err != nil {
		return nil, err
	}
	return &OrderProduser{
		conn: conn,
	}, nil
}
func (p *OrderProduser) SendAddOrderEvent(order Order) error {
	orderJson, err := json.Marshal(order)
	if err != nil {
		return err
	}
	p.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = p.conn.WriteMessages(kafka.Message{Value: orderJson})
	if err != nil {
		return err
	}
	return nil
}
