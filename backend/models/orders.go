package models

import (
	"fmt"

	"fyp.com/m/common"
	"fyp.com/m/db"
	"fyp.com/m/kafka_module"

	"log"
)

type Order struct {
	ID              int64  `json:"id"`
	BuyerID         int64  `json:"buyerId"`
	SellerID        int64  `json:"sellerId"`
	ProductID       int64  `json:"productId"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	ShippingAddress string `json:"shippingAddress"`
	Country         string `json:"country"`
	State           string `json:"state"`
	City            string `json:"city"`
	PostalCode      int64  `json:"postalCode"`
	PhoneNumber     int64  `json:"phoneNumber"`
	DeliveryOption  string `json:"deliveryOption"`
	CheckoutPrice   int64  `json:"checkoutPrice"`
	OrderStatus     string `json:"orderStatus"`
	PaymentMethod   string `json:"paymentMethod"`
	Time            int64  `json:"time"`
}

func GetAllOrders() ([]Order, error) {
	query := "SELECT * FROM orders"
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Printf("Error querying database: %v\n", err)
		return nil, err
	}
	defer rows.Close()
	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.BuyerID, &order.SellerID, &order.ProductID, &order.Name, &order.Email, &order.ShippingAddress, &order.Country,
			&order.State, &order.City, &order.PostalCode, &order.PhoneNumber, &order.DeliveryOption, &order.CheckoutPrice, &order.OrderStatus, &order.PaymentMethod, &order.Time)

		if err != nil {
			log.Printf("Error scanning rows: %v\n", err)
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (o *Order) Save() error {
	query := `INSERT INTO orders(
		buyer_id, seller_id, product_id, name, email, shipping_address, 
		country, state, city, postal_code, phone_number, delivery_option, 
		checkout_price, order_status, payment_method, time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		log.Printf("Error preparing save order query: %v\n", err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(o.BuyerID, o.SellerID, o.ProductID, o.Name, o.Email, o.ShippingAddress, o.Country, o.State, o.City, o.PostalCode,
		o.PhoneNumber, o.DeliveryOption, o.CheckoutPrice, o.OrderStatus, o.PaymentMethod, o.Time)

	if err != nil {
		log.Printf("Error executing save order query: %v\n", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error retrieving last insert ID: %v\n", err)
		return err
	}
	o.ID = id
	return nil

}

func GetOrderByID(id int64) (*Order, error) {
	query := "SELECT * FROM orders WHERE id = ?"
	row := db.DB.QueryRow(query, id)
	var order Order
	err := row.Scan(
		&order.ID, &order.BuyerID, &order.SellerID, &order.ProductID, &order.Name,
		&order.Email, &order.ShippingAddress, &order.Country, &order.State,
		&order.City, &order.PostalCode, &order.PhoneNumber, &order.DeliveryOption,
		&order.CheckoutPrice, &order.OrderStatus, &order.PaymentMethod, &order.Time,
	)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func GetOrderByBuyerID(id int64) ([]Order, error) {
	query := "SELECT * FROM orders WHERE buyer_id = ?"
	rows, err := db.DB.Query(query, id)
	if err != nil {
		log.Printf("Error querying database: %v\n", err)
		return nil, err
	}
	defer rows.Close()
	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.BuyerID, &order.SellerID, &order.ProductID, &order.Name, &order.Email, &order.ShippingAddress, &order.Country,
			&order.State, &order.City, &order.PostalCode, &order.PhoneNumber, &order.DeliveryOption, &order.CheckoutPrice, &order.OrderStatus, &order.PaymentMethod, &order.Time)

		if err != nil {
			log.Printf("Error scanning rows: %v\n", err)
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func GetOrderBySellerID(id int64) ([]Order, error) {
	query := "SELECT * FROM orders WHERE seller_id = ?"
	rows, err := db.DB.Query(query, id)
	if err != nil {
		log.Printf("Error querying database: %v\n", err)
		return nil, err
	}
	defer rows.Close()
	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.BuyerID, &order.SellerID, &order.ProductID, &order.Name, &order.Email, &order.ShippingAddress, &order.Country,
			&order.State, &order.City, &order.PostalCode, &order.PhoneNumber, &order.DeliveryOption, &order.CheckoutPrice, &order.OrderStatus, &order.PaymentMethod, &order.Time)

		if err != nil {
			log.Printf("Error scanning rows: %v\n", err)
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (order *Order) UpdateOrderStatus() error {
	query := `UPDATE orders SET order_status = ? WHERE id = ?`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("error while preparing query to update order status: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(order.OrderStatus, order.ID)
	if err != nil {
		return fmt.Errorf("error while executing query to update order status: %v", err)
	}
	user, err := GetUserbyID(order.BuyerID)
	if err != nil {
		return fmt.Errorf("error while getting user for sending email %v", err)
	}
	product, err := GetProductByID(order.ProductID)
	if err != nil {
		return fmt.Errorf("error while getting product for sending email %v", err)
	}
	msg := common.EmailMessage{
		Email:   user.Email,
		Subject: fmt.Sprintf("Update on Order #%d", order.ID),
		Body:    fmt.Sprintf("Your status for the order of the Product %s has been updated to %s", product.Name, order.OrderStatus),
	}
	go kafka_module.Producer(msg)
	return nil
}

func (order Order) DeleteOrder() error {
	query := "DELETE FROM orders WHERE id = ?"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(order.ID)
	return err
}
