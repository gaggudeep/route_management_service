package service

import (
	"errors"
	"fmt"
	"log"
	"microservice/src/entity"
	"time"
)

type (
	NaiveRouteCalcStrategy struct {
		orderGraphSvc OrderGraphSvc
	}
)

func NewNaiveRouteCalcStrategy(orderGraphSvc OrderGraphSvc) RouteCalculationStrategy {
	return &NaiveRouteCalcStrategy{
		orderGraphSvc: orderGraphSvc,
	}
}

func (s NaiveRouteCalcStrategy) CalculateRoute(agent entity.DeliveryAgent,
	orders []entity.Order) (entity.Route, error) {
	graph, err := s.orderGraphSvc.ConstructGraph(agent, orders)
	consumersMap := constructConsumersMap(orders)

	if err != nil {
		log.Printf("[ERROR] error constructing graph: %s", err)
		return entity.Route{}, err
	}
	if len(graph.Nodes) > 20 {
		err = errors.New("unsupported number of nodes for route calculation")
		log.Printf("[ERROR] error constructing graph: %s", err.Error())
		return entity.Route{}, err
	}

	routeNodes := make([]entity.RouteNode, 0, len(graph.Nodes))
	minCost := time.Hour * 24
	curCost, _ := time.ParseDuration("0h")
	curRoutesNode := make([]entity.RouteNode, 0, len(graph.Nodes))

	calculateOptimalRoute(0, curCost, graph, &minCost, &routeNodes, &curRoutesNode, consumersMap)

	return entity.Route{
		RouteNodes: routeNodes,
	}, nil
}

func calculateOptimalRoute(srcId int, curCost time.Duration, graph entity.Graph[time.Duration], minCost *time.Duration,
	bestRoute *[]entity.RouteNode, curRoute *[]entity.RouteNode, consumersMap map[int]map[int]bool) {
	node := graph.Nodes[srcId]

	node.VisitStatus = entity.VisitStatusVisited
	graph.Nodes[srcId] = node
	*curRoute = append(*curRoute, entity.RouteNode{
		MinutesToReachFromStart: curCost.Minutes(),
		OriginalNode:            node.OriginalNode,
	})

	if len(*curRoute) == len(graph.Nodes) {
		if curCost.Minutes() < minCost.Minutes() {
			*minCost = curCost
			newBestRoute := make([]entity.RouteNode, len(*curRoute))
			copy(newBestRoute, *curRoute)
			*bestRoute = newBestRoute
		}
	}
	if node.OriginalNode.GetType() == entity.NodeTypeRestaurant {
		markConsumerNodes(srcId, graph, consumersMap, entity.VisitStatusUnVisited)
		defer markConsumerNodes(srcId, graph, consumersMap, entity.VisitStatusUnVisitable)
	}
	for _, dest := range graph.Adj[srcId] {
		newCost, _ := time.ParseDuration(fmt.Sprintf("%fm", curCost.Minutes()+dest.Cost.Minutes()))
		nextNode := graph.Nodes[dest.DestinationNodeId]

		if newCost.Minutes() >= minCost.Minutes() ||
			nextNode.VisitStatus == entity.VisitStatusVisited ||
			nextNode.VisitStatus == entity.VisitStatusUnVisitable {
			continue
		}
		calculateOptimalRoute(dest.DestinationNodeId, newCost, graph, minCost, bestRoute, curRoute, consumersMap)
	}
	node.VisitStatus = entity.VisitStatusUnVisited
	graph.Nodes[srcId] = node
	*curRoute = (*curRoute)[:len(*curRoute)-1]
}

func markConsumerNodes(id int, graph entity.Graph[time.Duration], consumersMap map[int]map[int]bool,
	newStatus entity.VisitStatus) {
	restaurantConsumers := consumersMap[graph.Nodes[id].OriginalNode.GetId()]

	for _, neighbor := range graph.Adj[id] {
		node := graph.Nodes[neighbor.DestinationNodeId]

		if node.OriginalNode.GetType() == entity.NodeTypeConsumer && restaurantConsumers[node.OriginalNode.GetId()] {
			node.VisitStatus = newStatus
			graph.Nodes[neighbor.DestinationNodeId] = node
		}
	}
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
