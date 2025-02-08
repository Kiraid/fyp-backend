package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "../api.db") // Assign to the global DB
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

func CreateTable() {


createMessageTable := `CREATE TABLE IF NOT EXISTS messages (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	senderid INTEGER NOT NULL,
	recieverid INTEGER NOT NULL,
	content TEXT NOT NULL,
	time INTEGER NOT NULL

);`
//FOREIGN KEY (user1id) REFERENCES users(id) ON DELETE CASCADE,
//FOREIGN KEY (user2id) REFERENCES users(id) ON DELETE CASCADE,
//FOREIGN KEY (senderid) REFERENCES users(id) ON DELETE CASCADE
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