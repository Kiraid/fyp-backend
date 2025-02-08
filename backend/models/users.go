package models

import (
	"errors"
	"fmt"
	"log"

	"fyp.com/m/db"
	"fyp.com/m/utils"
)

type User struct{
	ID int64
	Email string `binding:"required"`
	Password string `binding:"required"`
	Username string 
	Role string 
}

func GetAllUsers() ([]User, error){
	query := "SELECT * FROM users"
	rows, err := db.DB.Query(query)
	if err != nil{
		log.Printf("Error querying database: %v\n", err)
	}
	defer rows.Close()
	var users []User
	for rows.Next(){
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.Username, &user.Role)
		if err != nil{
			log.Printf("Error scanning rows: %v\n", err)
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *User) Save() error {
	query := "INSERT INTO users(email, password, username, role) VALUES(?,?,?,?)"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		log.Printf("Error preparing while saving user query: %v\n", err)
	}
	hashedpass, err := utils.HashPassword(u.Password)
	if err != nil {
		log.Printf("Error hashing the pasword %v\n", err)
	}
	defer stmt.Close()
	result, err := stmt.Exec(u.Email, hashedpass, u.Username, u.Role)
	if err != nil {
		log.Printf("Error executing save user query: %v\n", err)
	}
	id, err := result.LastInsertId()
	u.ID = id
	return err
}

func Delete(id int64) error {
	query := "DELETE FROM users WHERE id = ?"
	result, err := db.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error while deleting user with id %d: %w", id, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error fetching rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with id %d", id)
	}
	return nil	
}

func GetUserbyID(id int64) (*User, error) {
	query := "SELECT * FROM users WHERE id = ?"
	row := db.DB.QueryRow(query, id)
	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Username, &user.Role)
	if err != nil{
		return nil, fmt.Errorf("error while getting user with id %d: %w", id, err)
	}
	return &user, nil
}

func (user User) UpdateUser() error {
	query := `UPDATE users SET email = ?, password = ?, username = ?, role = ? WHERE id = ?`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("error while preparing query for update user")
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Email, user.Password, user.Username, user.Role, user.ID)
	return err

}

func (u *User) ValidateCredential() error {
	query := "SELECT id, username, password, role FROM users WHERE email = ?"
	row := db.DB.QueryRow(query, u.Email)

	var retrievedPassword string
	err := row.Scan(&u.ID, &u.Username, &retrievedPassword, &u.Role)
	if err != nil {
		log.Printf("error while retrieving password from database")
	}
	passwordisValid := utils.CheckPassword(u.Password, retrievedPassword)
	if !passwordisValid{
		return errors.New("credentail invalid")
	}
	return nil

}