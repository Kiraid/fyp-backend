package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var DB *sql.DB

func InitDB() {
	var err error

	// MySQL connection string: "user:password@tcp(host:port)/database"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("MYSQL_USER"), // Add this if needed, default is "root"
		os.Getenv("MYSQL_ROOT_PASSWORD"),
		os.Getenv("MYSQL_HOST_NAME"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DB"),
	)

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

func CreateTable() {

	createMessageTable := `CREATE TABLE IF NOT EXISTS messages (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	senderid INTEGER NOT NULL,
	recieverid INTEGER NOT NULL,
	content TEXT NOT NULL,
	time INTEGER NOT NULL

);`
	// FOREIGN KEY (user1id) REFERENCES users(id) ON DELETE CASCADE,
	// FOREIGN KEY (user2id) REFERENCES users(id) ON DELETE CASCADE,
	// FOREIGN KEY (senderid) REFERENCES users(id) ON DELETE CASCADE
	_, err := DB.Exec(createMessageTable)
	if err != nil {
		log.Fatalf("Error creating message table: %v", err)
	}

	// Create indexes for faster lookups
	createIndexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_recieverid ON messages (recieverid);",
		"CREATE INDEX IF NOT EXISTS idx_senderid ON messages (senderid);",
	}

	for _, query := range createIndexes {
		_, err = DB.Exec(query)
		if err != nil {
			log.Fatalf("Error creating index: %v", err)
		}
	}
}
