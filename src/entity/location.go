package entity

type LocationDetails interface {
	GetId() int
	GetType() NodeType
}

type Location struct {
	Id             int            `json:"id"`
	GeoCoordinates GeoCoordinates `json:"geo_coordinates"`
	Type           NodeType       `json:"type"`
}

type GeoCoordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
