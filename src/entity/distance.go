package entity

type DistanceUnit int

const (
	DistanceUnitKm DistanceUnit = iota
)

type Distance struct {
	Value float64
	Unit  DistanceUnit
}
