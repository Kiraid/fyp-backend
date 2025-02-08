package models

import (
	"log"
	"fyp.com/m/db"
)

type Review struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"userId"`
	UserName  string `json:"username"`
	ProductID int64  `json:"productId"`
	Rating    int    `json:"rating"`
	Review    string `json:"review"`
}

func GetAllReviews() ([]Review, error) {
	query := "SELECT * FROM reviews"
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Printf("Error querying database: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var reviews []Review
	for rows.Next() {
		var review Review
		err := rows.Scan(&review.ID, &review.UserID, &review.UserName, &review.ProductID, &review.Rating, &review.Review)
		if err != nil {
			log.Printf("Error scanning rows: %v\n", err)
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

func (r *Review) Save() error {
	query := "INSERT INTO reviews(userid, username, productid, rating, review) VALUES(?,?,?,?,?)"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		log.Printf("Error preparing while saving review query: %v\n", err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(r.UserID, r.UserName, r.ProductID, r.Rating, r.Review)
	if err != nil {
		log.Printf("Error executing save review query: %v\n", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error fetching last insert ID for review: %v\n", err)
		return err
	}

	r.ID = id
	return nil
}

func GetReviewByID(id int64) (*Review, error) {
	query := "SELECT * FROM reviews WHERE id = ?"
	row := db.DB.QueryRow(query, id)

	var review Review
	err := row.Scan(&review.ID, &review.UserID, &review.UserName, &review.ProductID, &review.Rating, &review.Review)
	if err != nil {
		log.Printf("Error fetching review by ID: %v\n", err)
		return nil, err
	}

	return &review, nil
}

func GetReviewsByProductID(productID int64) ([]Review, error) {
	query := "SELECT * FROM reviews WHERE productid = ?"
	rows, err := db.DB.Query(query, productID)
	if err != nil {
		log.Printf("Error querying database for reviews: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var reviews []Review
	for rows.Next() {
		var review Review
		err := rows.Scan(&review.ID, &review.UserID, &review.UserName, &review.ProductID, &review.Rating, &review.Review)
		if err != nil {
			log.Printf("Error scanning rows for reviews: %v\n", err)
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}
