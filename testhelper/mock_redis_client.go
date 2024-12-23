package testhelper

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
)

func NewMockRedis(t *testing.T) *redis.Client {
	t.Helper()

	// redisサーバを作る
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("unexpected error while createing test redis server '%#v'", err)
	}
	// *redis.Clientを用意
	client := redis.NewClient(&redis.Options{
		Addr:     s.Addr(),
		Password: "",
		DB:       0,
	})
	return client
}
