package entity

import (
	"encoding/json"
	"fmt"
)

type NodeType int

func (nodeType NodeType) MarshalJSON() ([]byte, error) {
	switch nodeType {
	case NodeTypeAgent:
		return json.Marshal("AGENT")
	case NodeTypeConsumer:
		return json.Marshal("CONSUMER")
	case NodeTypeRestaurant:
		return json.Marshal("RESTAURANT")
	default:
		return nil, fmt.Errorf(`"%d" is not a valid node type`, nodeType)
	}
}

const (
	NodeTypeRestaurant NodeType = iota
	NodeTypeConsumer
	NodeTypeAgent
)

type RouteNode struct {
	MinutesToReachFromStart float64         `json:"minutesToReachFromStart"`
	OriginalNode            LocationDetails `json:"originalNode"`
}

type Route struct {
	RouteNodes []RouteNode `json:"routeNodes"`
}
