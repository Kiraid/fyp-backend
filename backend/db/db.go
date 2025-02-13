package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" // Import SQLite driver
)

var DB *sql.DB

// InitDB initializes the MySQL database connection
func InitDB() {
	var err error

	// MySQL connection string: "user:password@tcp(host:port)/database"
	dsn := "root:password@tcp(mysql:3306)/ecommerce_db?parseTime=true"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Could not establish a connection to the database:", err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
}

// CreateTable creates necessary tables in MySQL
func CreateTable() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users(
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		username VARCHAR(255) NOT NULL,
		role VARCHAR(50) NOT NULL DEFAULT 'buyer'
	)
	`
	_, err := DB.Exec(createUsersTable)
	if err != nil {
		log.Fatal("Could not create users table:", err)
	}

	createSellerDescriptionTable := `
	CREATE TABLE IF NOT EXISTS seller_description(
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		about TEXT NOT NULL,
		product_type VARCHAR(255) NOT NULL,
		user_id BIGINT NOT NULL UNIQUE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)
	`
	_, err = DB.Exec(createSellerDescriptionTable)
	if err != nil {
		log.Fatal("Could not create seller_description table:", err)
	}

	createProductCategoriesTable := `
	CREATE TABLE IF NOT EXISTS productscategories(
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		imagepath VARCHAR(255) NOT NULL,
		user_id BIGINT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)
	`
	_, err = DB.Exec(createProductCategoriesTable)
	if err != nil {
		log.Fatal("Could not create products categories table:", err)
	}

	createProductsTable := `
	CREATE TABLE IF NOT EXISTS products(
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		imagepath VARCHAR(255) NOT NULL,
		user_id BIGINT NOT NULL,
		category_name VARCHAR(255) NOT NULL,
		price DECIMAL(10,2) NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)
	`
	_, err = DB.Exec(createProductsTable)
	if err != nil {
		log.Fatal("Could not create products table:", err)
	}

	createOrdersTable := `
	CREATE TABLE IF NOT EXISTS orders(
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		buyer_id BIGINT NOT NULL,
		seller_id BIGINT NOT NULL,
		product_id BIGINT NOT NULL,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		shipping_address TEXT NOT NULL,
		country VARCHAR(100) NOT NULL,
		state VARCHAR(100) NOT NULL,
		city VARCHAR(100) NOT NULL,
		postal_code VARCHAR(20) NOT NULL,
		phone_number VARCHAR(20) NOT NULL,
		delivery_option VARCHAR(100) NOT NULL,
		checkout_price DECIMAL(10,2) NOT NULL,
		order_status VARCHAR(50) NOT NULL,
		payment_method VARCHAR(50) NOT NULL,
		time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (buyer_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (seller_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
	)
	`
	_, err = DB.Exec(createOrdersTable)
	if err != nil {
		log.Fatal("Could not create orders table:", err)
	}

	createReviewsTable := `
	CREATE TABLE IF NOT EXISTS reviews(
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		userid BIGINT NOT NULL,
		username VARCHAR(255) NOT NULL,
		productid BIGINT NOT NULL,
		rating INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
		review TEXT NOT NULL,
		FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (productid) REFERENCES products(id) ON DELETE CASCADE
	)
	`
	_, err = DB.Exec(createReviewsTable)
	if err != nil {
		log.Fatal("Could not create reviews table:", err)
	}
}
