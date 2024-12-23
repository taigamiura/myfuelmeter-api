package storage

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// グローバル変数
var Rdb *redis.Client // グローバル変数として宣言

// Redisの初期化
func InitRedis(redisAddr string, redisPassword string) {
	if redisAddr == "" {
		log.Fatal("Redis address is not set in the configuration.")
	}
	if redisPassword == "" {
		log.Println("Warning: Redis password is not set. Proceeding without authentication.")
	}

	Rdb = redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		Password:     redisPassword,
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	for i := 0; i < 5; i++ {
		_, err := Rdb.Ping(context.Background()).Result()
		if err == nil {
			log.Println("Connected to Redis")
			return
		}
		log.Printf("Failed to connect to Redis (attempt %d): %v\n", i+1, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatal("Failed to connect to Redis after multiple attempts.")
}
