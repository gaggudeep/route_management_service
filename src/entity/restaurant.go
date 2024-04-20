package entity

import "time"

type Restaurant struct {
	Location
	MeanPreparationTime time.Duration
}

func (r Restaurant) GetType() NodeType {
	return NodeTypeRestaurant
}

func (r Restaurant) GetId() int {
	return r.Id
}
