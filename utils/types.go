package utils

import (
	"time"
)

// 座標の構造体
type Coordinates struct {
	Message             string  `json:"message" validate:"required"`
	Timestamp           string  `json:"timestamp" validate:"required"`
	Latitude            float64 `json:"latitude" validate:"required,numeric"`
	Longitude           float64 `json:"longitude" validate:"required,numeric"`
	TraceID             string  `json:"traceId"`
	ProcessingStartTime string  `json:"processingStarttime"`
}

// トリップの構造体
type Trip struct {
	ID        uint      `gorm:"primaryKey"`
	StartLat  float64   `gorm:"column:start_latitude"`
	StartLong float64   `gorm:"column:start_longitude"`
	EndLat    float64   `gorm:"column:end_latitude"`
	EndLong   float64   `gorm:"column:end_longitude"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Track struct {
	Timestamp string  `json:"timestamp" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required,numeric"`// 緯度
	Longitude float64 `json:"longitude" validate:"required,numeric"`// 経度
}

// メッセージの構造体
type ResponseMessage struct {
	Status              string `json:"status"`
	ProcessingStartTime string `json:"processingStartTime"`
	TraceID             string `json:"traceID"`
}

type TrackFinishResponse struct {
	Distance  float64 `json:"distance" validate:"required,numeric"`
	FuelCost  float64 `json:"fuelCost" validate:"required,numeric"`
	StartTime string  `json:"startTime" validate:"required"`
	EndTime   string  `json:"endTime" validate:"required"`
}
