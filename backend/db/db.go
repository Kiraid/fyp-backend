package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./api.db") // Assign to the global DB
	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}
	_, err = DB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Could not enable foreign key support:", err)
	}
	
	err = DB.Ping()
	if err != nil {
		log.Fatal("Could not establish a connection to the database:", err)
	}
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
}

// CreateTable creates necessary tables
func CreateTable() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		username TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'buyer'
	)
	`
	_, err := DB.Exec(createUsersTable)
	if err != nil {
		log.Fatal("Could not create users table:", err)
	}

	createSellerDescrptionTable := `
	CREATE TABLE IF NOT EXISTS seller_description(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		about TEXT NOT NULL,
		product_type TEXT NOT NULL,
		user_id INTEGER NOT NULL UNIQUE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)
	`
	_, err = DB.Exec(createSellerDescrptionTable)
    if err != nil {
        log.Fatal("Could not create seller_description table:", err)
    }

	createProductCategoriesTable := `
	CREATE TABLE IF NOT EXISTS productscategories(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    imagepath TEXT NOT NULL,
	user_id INTEGER NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)
`
_, err = DB.Exec(createProductCategoriesTable)
if err != nil {
    log.Fatal("Could not create products table:", err)
}

	createProductsTable := `
CREATE TABLE IF NOT EXISTS products(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    imagepath TEXT NOT NULL,
	user_id INTEGER NOT NULL,
	category_name TEXT NOT NULL,
	price INTEGER NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)
`
_, err = DB.Exec(createProductsTable)
if err != nil {
    log.Fatal("Could not create products table:", err)
}

createOrdersTable := `
CREATE TABLE IF NOT EXISTS orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    buyer_id INTEGER NOT NULL,
    seller_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    shipping_address TEXT NOT NULL,
    country TEXT NOT NULL,
    state TEXT NOT NULL,
    city TEXT NOT NULL,
    postal_code INTEGER NOT NULL,
    phone_number INTEGER NOT NULL,
    delivery_option TEXT NOT NULL,
    checkout_price INTEGER NOT NULL,
    order_status TEXT NOT NULL,
    payment_method TEXT NOT NULL,
	time INTEGER NOT NULL,
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
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    userid INTEGER NOT NULL,
    username TEXT NOT NULL,
    productid INTEGER NOT NULL,
    rating INTEGER NOT NULL,
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
