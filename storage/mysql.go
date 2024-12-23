package storage

import (
	"log"

	"github.com/taigamiura/my-fuel-meter-project/api/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

// データベースの初期化
func InitDB(databaseDsn string) {
	var err error
	Db, err = gorm.Open(mysql.Open(databaseDsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Connected to database")

	if err := Db.AutoMigrate(&utils.Trip{}); err != nil {
		log.Println("Error during auto migration:", err)
	} else {
		log.Println("Database migrated successfully")
	}
}
