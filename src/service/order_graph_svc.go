package service

import (
	"log"
	"microservice/src/entity"
	"time"
)

type (
	OrderGraphSvcImpl struct {
		travelTimeSvc TravelTimeSvc
		distanceSvc   DistanceSvc
	}
)

func NewOrderGraphSvvImpl(travelTimeSvc TravelTimeSvc, distanceTimeSvc DistanceSvc) OrderGraphSvc {
	return &OrderGraphSvcImpl{
		travelTimeSvc: travelTimeSvc,
		distanceSvc:   distanceTimeSvc,
	}
}

func (s OrderGraphSvcImpl) ConstructGraph(agent entity.DeliveryAgent,
	orders []entity.Order) (entity.Graph[time.Duration], error) {
	nodes := constructNodes(agent, orders)
	edges, err := s.constructEdges(agent, nodes, orders)

	if err != nil {
		log.Printf("[ERROR] error from travel time svc: %s", err.Error())
		return entity.Graph[time.Duration]{}, nil
	}

	return entity.Graph[time.Duration]{
		Nodes: nodes,
		Adj:   edges,
	}, nil
}

func (s OrderGraphSvcImpl) constructEdges(agent entity.DeliveryAgent, nodes []entity.Node[time.Duration],
	orders []entity.Order) ([][]entity.Edge[time.Duration], error) {
	adj := make([][]entity.Edge[time.Duration], 0, len(nodes))
	consumersMap := constructConsumersMap(orders)

	for i := 0; i < len(nodes); i++ {
		src := nodes[i]

		adj = append(adj, []entity.Edge[time.Duration]{})
		for j := 0; j < len(nodes); j++ {
			dest := nodes[j]

			if isConnected(src, dest, consumersMap) {
				dist := s.distanceSvc.CalculateDistance(src.GeoCoordinates, dest.GeoCoordinates)
				travelTime, err := s.travelTimeSvc.CalculateTravelTime(agent, dist)

				if err != nil {
					return nil, err
				}
				adj[i] = append(adj[i], entity.Edge[time.Duration]{
					DestinationNodeId: j,
					Cost:              travelTime + dest.OverheadCost,
				})
			}
		}
	}

	return adj, nil
}

/*
isConnected any restaurant can be visited from any node (agent, other restaurant, consumer)
but a consumer can only be visited from a particular restaurant
*/
func isConnected(src entity.Node[time.Duration], dest entity.Node[time.Duration],
	consumerMap map[int]map[int]bool) bool {
	if dest.OriginalNode.GetType() == entity.NodeTypeRestaurant {
		return true
	}

	consumersSet, ok := consumerMap[src.OriginalNode.GetId()]
	if !ok {
		return false
	}
	return consumersSet[dest.OriginalNode.GetId()]
}

func constructNodes(agent entity.DeliveryAgent, orders []entity.Order) []entity.Node[time.Duration] {
	// 2x denotes that each order corresponds to 1 restaurant and 1 consumer
	// + 1 to add a  node for delivery agent's initial location
	nodes := make([]entity.Node[time.Duration], 0, 2*len(orders)+1)

	nodes = append(nodes, entity.Node[time.Duration]{
		Id:             0,
		OriginalNode:   agent,
		VisitStatus:    entity.VisitStatusUnVisited,
		GeoCoordinates: agent.GeoCoordinates,
	})
	for i, order := range orders {
		nodes = append(nodes, entity.Node[time.Duration]{
			Id:             i + 1,
			OriginalNode:   order.Restaurant,
			VisitStatus:    entity.VisitStatusUnVisited,
			GeoCoordinates: order.Restaurant.GeoCoordinates,
			OverheadCost:   order.Restaurant.MeanPreparationTime,
		})
		nodes = append(nodes, entity.Node[time.Duration]{
			Id:             i + 2,
			OriginalNode:   order.Consumer,
			VisitStatus:    entity.VisitStatusUnVisited,
			GeoCoordinates: order.Consumer.GeoCoordinates,
		})
	}

	return nodes
}

func constructConsumersMap(orders []entity.Order) map[int]map[int]bool {
	consumerMap := make(map[int]map[int]bool)

	for _, order := range orders {
		consumersSet, ok := consumerMap[order.Restaurant.Id]
		if !ok {
			consumersSet = make(map[int]bool)
		}
		consumersSet[order.Consumer.Id] = true
		consumerMap[order.Restaurant.Id] = consumersSet
	}

	return consumerMap
}
