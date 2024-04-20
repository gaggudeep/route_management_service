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
	edges, err := s.constructEdges(agent, nodes)

	if err != nil {
		log.Printf("[ERROR] error from travel time svc: %s", err.Error())
		return entity.Graph[time.Duration]{}, nil
	}

	return entity.Graph[time.Duration]{
		Nodes: nodes,
		Adj:   edges,
	}, nil
}

func (s OrderGraphSvcImpl) constructEdges(agent entity.DeliveryAgent,
	nodes []entity.Node[time.Duration]) ([][]entity.Edge[time.Duration], error) {
	adj := make([][]entity.Edge[time.Duration], 0, len(nodes))

	for i := 0; i < len(nodes); i++ {
		src := nodes[i]

		adj = append(adj, []entity.Edge[time.Duration]{})
		for j := 0; j < len(nodes); j++ {
			dest := nodes[j]
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

	return adj, nil
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
			VisitStatus:    entity.VisitStatusUnVisitable,
			GeoCoordinates: order.Consumer.GeoCoordinates,
		})
	}

	return nodes
}
