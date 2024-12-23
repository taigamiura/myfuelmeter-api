package utils

import (
	"time"

	"github.com/google/uuid"
)

// ユニークなIDを生成するヘルパー関数
func GenerateUniqueID() string {
	id := uuid.New().String()
	return id[:29]
}

// フォーマット時間関数
func FormatTime(input string) (string, error) {
	layout := time.RFC3339 // "2006-01-02T15:04:05Z"

	t, err := time.Parse(layout, input)
	if err != nil {
		return "", err
	}

	return t.Format("20060102150405"), nil
}
