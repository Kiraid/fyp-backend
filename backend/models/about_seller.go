package models

import (
	"fyp.com/m/db"
	"log"
)

type Seller_about struct {
	ID          int64  `json:"id"`
	About       string `json:"about" binding:"required"`
	ProductType string `json:"product_type" `
	UserID      int64  `json:"userId" `
}

func (u *Seller_about) Save() error {
	query := "INSERT INTO seller_description(about, product_type, user_id) VALUES(?,?,?)"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		log.Printf("Error preparing while saving seller_description query: %v\n", err)
	}

	defer stmt.Close()
	result, err := stmt.Exec(u.About, u.ProductType, u.UserID)

	if err != nil {
		log.Printf("Error executing save seller_description query: %v\n", err)
	}
	id, err := result.LastInsertId()
	u.ID = id
	return err
}
