package models

import (
	"log"
	"fyp.com/m/db"
)

type Product struct{
	ID int64 `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	ImagePath string `json:"imagepath"`
	UserID int64 `json:"userId"`
	Category_name string `json:"categoryName"`
	Price int64 `json:"price"`
}

func GetAllProducts() ([]Product, error){
	query := "SELECT * FROM products"
	rows, err := db.DB.Query(query)
	if err != nil{
		log.Printf("Error querying database: %v\n", err)
	}
	defer rows.Close()
	var products []Product
	for rows.Next(){
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.ImagePath, &product.UserID, &product.Category_name, &product.Price)
		if err != nil{
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
	}

	defer stmt.Close()
	result, err := stmt.Exec(u.Name, u.Description, u.ImagePath, u.UserID, u.Category_name, u.Price)
	
	if err != nil {
		log.Printf("Error executing save product query: %v\n", err)
	}
	id, err := result.LastInsertId()
	u.ID = id
	return err
}

func GetProductByID(id int64) (*Product, error){
	query := "SELECT * FROM products WHERE id = ?"
	row := db.DB.QueryRow(query, id)
	var product Product
	err := row.Scan(&product.ID, &product.Name, &product.Description, &product.ImagePath, &product.UserID, &product.Category_name, &product.Price)
	if err != nil {
		return nil, err
	}
	return &product,   nil
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
		params = append(params, product.UserID )
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

func (product Product) DeleteProduct() error{
	query := "DELETE FROM products WHERE id = ?"
	stmt,err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(product.ID)
	return err
}

func GetProductsbyUserID(id int64) ([]Product, error){
	query := "SELECT * FROM products where user_id = ?"
	rows, err := db.DB.Query(query, id)
	if err != nil{
		log.Printf("Error querying database: %v\n", err)
	}
	defer rows.Close()
	var products []Product
	for rows.Next(){
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.ImagePath, &product.UserID, &product.Category_name, &product.Price)
		if err != nil{
			log.Printf("Error scanning rows: %v\n", err)
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}


func GetProductsbySearch(search string) ([]Product, error){
	query := "SELECT * FROM products WHERE name LIKE ?"
	searchTerm := "%" + search + "%" 
	rows, err := db.DB.Query(query, searchTerm)
	if err != nil {
		log.Printf("Error querying database for searchbar: %v\n", err)
	}
	defer rows.Close()
	var products []Product
	for rows.Next(){
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.ImagePath, &product.UserID, &product.Category_name, &product.Price)
		if err != nil{
			log.Printf("Error scanning rows for searchbar: %v\n", err)
			return nil, err
		}
		products = append(products, product)
	}
	return products,nil
}


