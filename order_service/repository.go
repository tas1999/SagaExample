package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type OrderRepository struct {
	db *sql.DB
}
type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslMode"`
}

func New(conf PostgresConfig) (*OrderRepository, error) {
	connStr := fmt.Sprintf("user=%s host=%s port=%d password=%s dbname=%s sslmode=%s",
		conf.Username, conf.Host, conf.Port, conf.Password, conf.DBName, conf.SSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &OrderRepository{
		db: db,
	}, nil
}
func (r *OrderRepository) AddOrder(order Order) (int, error) {
	row := r.db.QueryRow("INSERT INTO orders (name,price,status,user_id) VALUES ($1,$2,$3,$4) returning id",
		order.Name, order.Price, order.Status, order.UserId)
	err := row.Err()
	if err != nil {
		return 0, err
	}
	err = row.Scan(&order.Id)
	if err != nil {
		return 0, err
	}
	return order.Id, nil
}
func (r *OrderRepository) ChangeOrder(order Order) error {
	row := r.db.QueryRow("UPDATE orders SET name=$1,price=$2,status=$3,user_id=$4 where id = $5",
		order.Name, order.Price, order.Status, order.UserId, order.Id)
	err := row.Err()
	if err != nil {
		return err
	}
	return nil
}
func (r *OrderRepository) GetOrder(orderId int) (*Order, error) {
	rows, err := r.db.Query("select id, name,price,status,user_id from orders where id=$1", orderId)
	if err != nil {
		return nil, err
	}
	o := Order{}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&o.Id, &o.Name, &o.Price, &o.Status, &o.UserId)
	return &o, err
}
