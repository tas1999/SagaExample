package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type CustomerRepository struct {
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

func New(conf PostgresConfig) (*CustomerRepository, error) {
	connStr := fmt.Sprintf("user=%s host=%s port=%d password=%s dbname=%s sslmode=%s",
		conf.Username, conf.Host, conf.Port, conf.Password, conf.DBName, conf.SSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &CustomerRepository{
		db: db,
	}, nil
}
func (r *CustomerRepository) AddCustomer(customer Customer) (int, error) {
	row := r.db.QueryRow("INSERT INTO customers (name,credit,user_id) VALUES ($1,$2,$3) returning id",
		customer.Name, customer.Credit, customer.UserId)
	err := row.Err()
	if err != nil {
		return 0, err
	}
	err = row.Scan(&customer.Id)
	if err != nil {
		return 0, err
	}
	return customer.Id, nil
}
func (r *CustomerRepository) ChangeCustomer(customer Customer) error {
	row := r.db.QueryRow("UPDATE customers SET name=$1,credit=$2,user_id=$3 where id = $4",
		customer.Name, customer.Credit, customer.UserId, customer.Id)
	err := row.Err()
	if err != nil {
		return err
	}
	return nil
}
func (r *CustomerRepository) GetCustomer(userId int) (*Customer, error) {
	rows, err := r.db.Query("select id, name,credit,user_id from customers where user_id=$1", userId)
	if err != nil {
		return nil, err
	}
	o := Customer{}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&o.Id, &o.Name, &o.Credit, &o.UserId)
	return &o, err
}
