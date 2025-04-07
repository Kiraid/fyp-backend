package storing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
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

var ES *elasticsearch.Client

func InitES() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200/",
		},
		APIKey: "",
	}
	var err error
	ES, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

}

func (p *Product) Save() error {
	jsonData, err := json.Marshal(p)
	log.Println(p)
	if err != nil {
		log.Printf("Error marshaling product data: %v\n", err)
		return err
	}

	req := bytes.NewReader(jsonData)
	res, err := ES.Index("products", req)
	if err != nil {
		log.Printf("Error indexing document in Elasticsearch: %v\n", err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		log.Printf("Unexpected response status from ES: %d\n", res.StatusCode)
		return fmt.Errorf("failed to store product in Elasticsearch")
	}

	log.Println("Product successfully stored in Elasticsearch")
	return nil
}
