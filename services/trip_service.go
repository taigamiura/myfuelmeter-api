package services

import (
	"github.com/taigamiura/my-fuel-meter-project/api/utils"
	"gorm.io/gorm"
)

// Create Trip 新規作成
func CreateTrip(db *gorm.DB, startCoords utils.Coordinates, endCoords utils.Coordinates) error {
	trip := utils.Trip{
		StartLat:  startCoords.Latitude,
		StartLong: startCoords.Longitude,
		EndLat:    endCoords.Latitude,
		EndLong:   endCoords.Longitude,
	}

	tx := db.Begin()
	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&trip).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetTripByID IDを取得
func GetTripByID(db *gorm.DB, id string) (*utils.Trip, error) {
	var Trip utils.Trip
	if err := db.First(&Trip, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &Trip, nil
}

// Update Trip 情報を更新
func UpdateTrip(db *gorm.DB, Trip *utils.Trip) error {
	return db.Save(Trip).Error
}

// Delete Trip削除
func DeleteTrip(db *gorm.DB, id string) error {
	if err := db.Delete(&utils.Trip{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
