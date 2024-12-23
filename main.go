package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-redis/redis/v8"
	"github.com/taigamiura/my-fuel-meter-project/api/config"
	"github.com/taigamiura/my-fuel-meter-project/api/future"
	"github.com/taigamiura/my-fuel-meter-project/api/storage"
	"gorm.io/gorm"
)

// グローバル変数
var (
	db  *gorm.DB
	rdb *redis.Client
	cfg *config.Config
)

// メイン関数
func main() {
	// 設定の読み込み
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	// DBとRedisの初期化
	storage.InitDB(cfg.DatabaseDsn)
	log.Println("MySQL initialized successfully.")

	storage.InitRedis(cfg.RedisAddr, cfg.RedisPassword)
	log.Println("Redis initialized successfully.")

	log.Printf("Server started on :%s", cfg.WebSocketPort)
	go handleShutdown()

	// WebSocketハンドラーの設定
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		future.HandleWebSocket(w, r, storage.Db, storage.Rdb) // Redisクライアントを渡す
	})

	// サーバーの開始
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.WebSocketPort), nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// シグナルハンドラーの設定
func handleShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Println("Shutting down gracefully...")
	if err := rdb.Close(); err != nil {
		log.Println("Error closing Redis connection:", err)
	} else {
		log.Println("Redis connection closed.")
	}
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Println("Error closing MySQL connection:", err)
			} else {
				log.Println("MySQL connection closed.")
			}
		}
	}
	log.Println("Connections closed.")
	os.Exit(0)
}
