package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func Product_key(id int) string {
	return fmt.Sprintf("product#%d", id)
}

var RDB *redis.Client

func InitRedis() {
	redisHost := os.Getenv("REDIS_HOST_NAME")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	RDB = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
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

//using msg pack instead of json because its faster

// func CacheProduct(prod Product) {
// 	ctx := context.Background()

// 	// Convert product struct to MessagePack
// 	productMsgPack, err := msgpack.Marshal(prod)
// 	if err != nil {
// 		log.Println("Error marshaling product:", err)
// 		return
// 	}

// 	key := Product_key(int(prod.ID))
// 	err = RDB.Set(ctx, key, productMsgPack, 0).Err() // No expiration here, set it in GetCachedProduct
// 	if err != nil {
// 		log.Println("Error caching product:", err)
// 	}
// }

// // GetCachedProduct retrieves the product from Redis and deserializes it using MessagePack
// func GetCachedProduct(id int) (*Product, error) {
// 	ctx := context.Background()
// 	key := Product_key(id)

// 	// Get MessagePack data from Redis
// 	productMsgPack, err := RDB.Get(ctx, key).Bytes() // Use Bytes() instead of Result() for binary data
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Convert MessagePack back to struct
// 	var product Product
// 	err = msgpack.Unmarshal(productMsgPack, &product)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Reset expiration time
// 	RDB.Expire(ctx, key, time.Hour)
// 	return &product, nil
// }
