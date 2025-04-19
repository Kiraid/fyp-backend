package common

//Common structs Imported in multiple packages for use
type EmailMessage struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type Product struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	ImagePath     string `json:"imagepath"`
	UserID        int64  `json:"userId"`
	Category_name string `json:"categoryName"`
	Price         int64  `json:"price"`
}
