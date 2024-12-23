package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/taigamiura/my-fuel-meter-project/api/utils"
)

// Redisに座標を保存する関数
func SaveCoordinatesToRedis(rdb *redis.Client, coords *utils.Coordinates, traceID string) error {
	key := "TRACE_ID:" + traceID
	track := utils.Track{
		Timestamp: coords.Timestamp,
		Latitude:  coords.Latitude,
		Longitude: coords.Longitude,
	}
	jsonData, err := json.Marshal(track)
	if err != nil {
		return fmt.Errorf("error encoding Coordinates to JSON: %w", err)
	}
	if rdb == nil {
		return fmt.Errorf("error Redis client is not initialized.%s: %w", key, err)
	}
	if err := rdb.RPush(context.Background(), key, jsonData).Err(); err != nil {
		return fmt.Errorf("error saving to Redis with key %s: %w", key, err)
	}

	return nil
}

// Redisから座標を一括取得する関数
func GetCoordinatesFromRedis(rdb *redis.Client, traceID string) ([]utils.Track, error) {
	key := "TRACE_ID:" + traceID

	if rdb == nil {
		return nil, fmt.Errorf("error Redis client is not initialized for key: %s", key)
	}

	// Redisからデータを一度に取得
	jsonData, err := rdb.LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("error retrieving data from Redis with key %s: %w", key, err)
	}

	// スライスの長さを予め設定（オプション）
	tracks := make([]utils.Track, 0, len(jsonData))

	// JSONデータのデコードを一括で行う
	for _, data := range jsonData {
		var track utils.Track
		if err := json.Unmarshal([]byte(data), &track); err != nil {
			return nil, fmt.Errorf("error decoding JSON data for key %s: %w", key, err)
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}
