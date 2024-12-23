package utils

import (
	"math"
)

// Position 構造体: 緯度・経度を格納
// type Position struct {
// 	Latitude  float64
// 	Longitude float64
// }

// haversineDistance 関数: 2 点間の距離を計算
func HaversineDistance(start, end Track) float64 {
	const EarthRadius = 6371.0 // 地球の半径 (km)

	// 度をラジアンに変換
	lat1 := degreesToRadians(start.Latitude)
	lon1 := degreesToRadians(start.Longitude)
	lat2 := degreesToRadians(end.Latitude)
	lon2 := degreesToRadians(end.Longitude)

	// 緯度と経度の差
	dLat := lat2 - lat1
	dLon := lon2 - lon1

	// Haversine 公式
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// 距離 (km)
	distance := EarthRadius * c
	return distance
}

// degreesToRadians 関数: 度をラジアンに変換
func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
