package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/taigamiura/my-fuel-meter-project/api/utils"
)

type TrackingStatus struct {
	TraceId string `json:"traceId"`
	TrackingStatus  string `json:"status"`
}

func SaveTrackingState(rdb *redis.Client, traceID, status string) error {
	key := "TRACKING_STATUS:" + traceID
	trackingStatus := TrackingStatus{
		TraceId: traceID,
		TrackingStatus:  status,
	}
	// JSONにエンコードしてRedisに保存
	jsonData, err := json.Marshal(trackingStatus)
	if err != nil {
		return fmt.Errorf("error encoding TrackingStatus to JSON: %w", err)
	}

	if rdb == nil {
		return fmt.Errorf("error Redis client is not initialized.%s: %w", key, err)
	}

	if err := rdb.Set(context.Background(), key, jsonData, 0).Err(); err != nil {
		return fmt.Errorf("error saving to Redis with key %s: %w", key, err)
	}

	return nil
}
func LoadTrackingState(rdb *redis.Client, traceID string) (*utils.Coordinates, error) {
	key := "TRACKING_STATUS:" + traceID
	val, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		return &utils.Coordinates{}, err
	}
	var trackingStatus map[string]interface{}
	if err := json.Unmarshal([]byte(val), &trackingStatus); err != nil {
		return &utils.Coordinates{}, err
	}
	// 現在の座標を取得
	coords := trackingStatus["coords"].(utils.Coordinates)
	return &coords, nil
}
