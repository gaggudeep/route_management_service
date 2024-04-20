package helper

import (
	"math"
	"microservice/src/entity"
)

// REFERENCES: https://en.wikipedia.org/wiki/Haversine_formula, https://en.wikipedia.org/wiki/Great-circle_distance
func CalculateCentralAngleRadians(origin entity.GeoCoordinates, destination entity.GeoCoordinates) float64 {
	longitudeAbsDiff := math.Abs(origin.Longitude - destination.Longitude)
	centralAngleRadians := math.Acos(
		math.Sin(origin.Latitude)*math.Sin(destination.Latitude) +
			math.Cos(origin.Latitude)*math.Cos(destination.Latitude)*math.Cos(longitudeAbsDiff),
	)

	return centralAngleRadians
}

func RoundTo(n float64, decimals int) float64 {
	return math.Round(n*math.Pow(10, float64(decimals))) / math.Pow(10, float64(decimals))
}
