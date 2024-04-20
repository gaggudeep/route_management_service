package entity

type TimeUnit int

const (
	TimeUnitHr TimeUnit = iota
)

type Speed struct {
	Value        float64      `json:"value"`
	DistanceUnit DistanceUnit `json:"distance_unit"`
	TimeUnit     TimeUnit     `json:"time_unit"`
}

type DeliveryAgent struct {
	Location
	AverageSpeed Speed `json:"average_speed"`
}

func (d DeliveryAgent) GetId() int {
	return d.Id
}

func (d DeliveryAgent) GetType() NodeType {
	return NodeTypeAgent
}
