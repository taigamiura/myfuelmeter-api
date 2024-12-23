package utils

import (
	"log"
	"strconv"

	"github.com/taigamiura/my-fuel-meter-project/api/config"
)

// 燃料費を計算する関数
func CalculateFuelCost(distance float64) float64 {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
		return 0
	}

	fuelPricePerLiter, err := strconv.ParseFloat(cfg.FuelPricePerLiter, 64)
	if err != nil {
		// デフォルト値を返す
		fuelPricePerLiter = 150.0
	}

	fuelEfficiency, err := strconv.ParseFloat(cfg.FuelEfficiency, 64)
	if err != nil {
		// デフォルト値を返す
		fuelEfficiency = 15.0
	}

	return (distance / fuelEfficiency) * fuelPricePerLiter
}
