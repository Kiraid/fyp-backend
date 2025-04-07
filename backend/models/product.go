package models

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"fyp.com/m/db"
)

type Product struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	ImagePath     string  `json:"imagepath"`
	UserID        int64   `json:"userId"`
	Category_name string  `json:"categoryName"`
	Price         float64 `json:"price"`
}

func GetAllProducts() ([]Product, error) {
	query := "SELECT * FROM products"
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Printf("Error querying database: %v\n", err)
	}
	defer rows.Close()
	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.ImagePath, &product.UserID, &product.Category_name, &product.Price)
		if err != nil {
			log.Printf("Error scanning rows: %v\n", err)
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (u *Product) Save() error {
	query := "INSERT INTO products(name, description, imagepath, user_id, category_name, price) VALUES(?,?,?,?,?,?)"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		log.Printf("Error preparing while saving product query: %v\n", err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(u.Name, u.Description, u.ImagePath, u.UserID, u.Category_name, u.Price)
	if err != nil {
		log.Printf("Error executing save product query: %v\n", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v\n", err)
		return err
	}
	u.ID = id

	// Run the HTTP request in a goroutine
	go func(product Product) {
		jsonData, err := json.Marshal(product)
		if err != nil {
			log.Printf("Error marshaling product data: %v\n", err)
			return
		}

		resp, err := http.Post("http://localhost:8083/save-ES", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Error sending request to Elasticsearch service: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Elasticsearch service responded with status: %d\n", resp.StatusCode)
		}
	}(*u) // Passing a copy of the product struct to avoid potential data race

	return nil
}

func GetProductByID(id int64) (*Product, error) {
	product, err := GetCachedProduct(int(id))
	if err == nil && product != nil {
		log.Printf("Loaded from redis")
		return product, nil
	}

	query := "SELECT * FROM products WHERE id = ?"
	row := db.DB.QueryRow(query, id)
	product = &Product{}
	err = row.Scan(&product.ID, &product.Name, &product.Description, &product.ImagePath, &product.UserID, &product.Category_name, &product.Price)
	if err != nil {
		return nil, err
	}
	CacheProduct(*product)
	log.Printf("Added to redis")
	return product, nil
}

func (product Product) UpdateProduct() error {
	query := `UPDATE products SET `
	params := []interface{}{}

	if product.Name != "" {
		query += "name = ?, "
		params = append(params, product.Name)
	}
	if product.Description != "" {
		query += "description = ?, "
		params = append(params, product.Description)
	}
	if product.ImagePath != "" {
		query += "imagepath = ?, "
		params = append(params, product.ImagePath)
	}
	if product.UserID != 0 {
		query += "user_id = ?, "
		params = append(params, product.UserID)
	}
	if product.Category_name != "" {
		query += "category_name = ?, "
		params = append(params, product.Category_name)
	}
	if product.Price != 0 {
		query += "price = ?, "
		params = append(params, product.Price)
	}

	query = query[:len(query)-2]
	query += " WHERE id = ?"
	params = append(params, product.ID)

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(params...)
	return err
}

func (product Product) DeleteProduct() error {
	query := "DELETE FROM products WHERE id = ?"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(product.ID)
	return err
}

func GetProductsbyUserID(id int64) ([]Product, error) {
	query := "SELECT * FROM products where user_id = ?"
	rows, err := db.DB.Query(query, id)
	if err != nil {
		log.Printf("Error querying database: %v\n", err)
	}
	defer rows.Close()
	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.ImagePath, &product.UserID, &product.Category_name, &product.Price)
		if err != nil {
			log.Printf("Error scanning rows: %v\n", err)
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func GetProductsbySearch(search string) ([]Product, error) {
	query := "SELECT * FROM products WHERE name LIKE ?"
	searchTerm := "%" + search + "%"
	rows, err := db.DB.Query(query, searchTerm)
	if err != nil {
		log.Printf("Error querying database for searchbar: %v\n", err)
	}
	defer rows.Close()
	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.ImagePath, &product.UserID, &product.Category_name, &product.Price)
		if err != nil {
			log.Printf("Error scanning rows for searchbar: %v\n", err)
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}
