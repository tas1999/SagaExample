package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Repository interface {
	AddOrder(order Order) (int, error)
	ChangeOrder(order Order) error
	GetOrder(orderId int) (*Order, error)
}
type Produser interface {
	SendAddOrderEvent(order Order) error
}
type Consumer interface {
	Start()
}

func main() {
	rep, err := New(PostgresConfig{
		Username: "postgres",
		Host:     "postgres",
		Port:     5432,
		Password: "postgres",
		DBName:   "order_service",
		SSLMode:  "disable",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	var prod *OrderProduser
	time.Sleep(time.Second * 20)
	prod, err = NewOrderProduser(OrderProduserConfig{
		Topic: "order.create",
		Host:  "kafka",
		Port:  9092,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	cons, err := NewConsumer(OrderProduserConfig{
		Topic: "order.check.result",
		Host:  "kafka",
		Port:  9092,
	}, rep, prod)
	if err != nil {
		fmt.Println(err)
		return
	}

	//var cons Consumer
	http.HandleFunc("/AddOrder", func(w http.ResponseWriter, r *http.Request) {
		var order Order
		defer r.Body.Close()
		d := json.NewDecoder(r.Body)
		err := d.Decode(&order)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		order.Status = Pending
		id, err := rep.AddOrder(order)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		order.Id = id
		err = prod.SendAddOrderEvent(order)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		en := json.NewEncoder(w)
		err = en.Encode(&order)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
	})
	go cons.Start()
	s := &http.Server{
		Addr: ":8080",
	}
	err = s.ListenAndServe()
	fmt.Println(err)
}

const (
	Default OrderStatus = iota
	Pending
	Approved
	Cancelled
)

type OrderStatus int
type Order struct {
	Id     int
	Name   string
	Price  int
	UserId int
	Status OrderStatus
}
