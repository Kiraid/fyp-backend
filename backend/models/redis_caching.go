package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func Product_key(id int) string {
	return fmt.Sprintf("product#%d", id)
}

var RDB *redis.Client

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "yourpassword",
		DB:       0,
	})

	ctx := context.Background()
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Could not connect to Redis:", err)
	}

	log.Println("Connected to Redis successfully!")
}

func CacheProduct(prod Product) {
	ctx := context.Background()

	// Convert product struct to JSON
	productJSON, err := json.Marshal(prod)

	if err != nil {
		log.Println("Error marshaling product:", err)
		return
	}

	key := Product_key(int(prod.ID))
	err = RDB.Set(ctx, key, productJSON, 0).Err()
	if err != nil {
		log.Println("Error caching product:", err)
	}
}

func GetCachedProduct(id int) (*Product, error) {
	ctx := context.Background()
	key := Product_key(id)
	// Get JSON from Redis
	productJSON, err := RDB.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	// Convert JSON back to struct
	var product Product
	err = json.Unmarshal([]byte(productJSON), &product)
	if err != nil {
		return nil, err
	}
	RDB.Expire(ctx, key, time.Hour)
	return &product, nil
}
