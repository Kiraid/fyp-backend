package models

import (
	"fmt"
	"log"

	"example.com/chat/db"
)

type Message struct {
	ID         int64 `json:"id"`
	SenderID   int64
	RecieverID int64  `json:"recieverId" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Time       int64  `json:"time"`
}

func (m *Message) CreateMessage() error {
	query := "INSERT INTO messages(senderid, recieverid, content, time) VALUES(?,?,?,?)"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		log.Printf("Error preparing while saving product query: %v\n", err)
	}

	defer stmt.Close()
	result, err := stmt.Exec(m.SenderID, m.RecieverID, m.Content, m.Time)

	if err != nil {
		log.Printf("Error executing save product query: %v\n", err)
	}
	id, err := result.LastInsertId()
	m.ID = id
	return err
}

func GetAllMessages(senderID, receiverID, limit, offset int) ([]Message, error) {
	query := `
		SELECT * FROM messages
		WHERE (senderid = ? AND recieverid = ?) OR (senderid = ? AND recieverid = ?)
		ORDER BY time ASC
		LIMIT ? OFFSET ?
	`
	rows, err := db.DB.Query(query, senderID, receiverID, receiverID, senderID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.RecieverID, &msg.Content, &msg.Time); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func GetUserNamebyId(id int64) (string, error) {
	query := "SELECT username FROM users WHERE id = ?"
	row := db.DB.QueryRow(query, id)
	var user string
	err := row.Scan(&user)
	if err != nil {
		return "", fmt.Errorf("error while getting user with id %d: %w", id, err)
	}
	return user, nil
}
