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

	calculateOptimalRoute(0, curCost, graph, &minCost, &routeNodes, &curRoutesNode)

	return entity.Route{
		RouteNodes: routeNodes,
	}, nil
}

func calculateOptimalRoute(src int, curCost time.Duration, graph entity.Graph[time.Duration], minCost *time.Duration,
	bestRoute *[]entity.RouteNode, curRoute *[]entity.RouteNode) {
	node := graph.Nodes[src]

	node.VisitStatus = entity.VisitStatusVisited
	graph.Nodes[src] = node
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
	for _, dest := range graph.Adj[src] {
		newCost, _ := time.ParseDuration(fmt.Sprintf("%fm", curCost.Minutes()+dest.Cost.Minutes()))

		if newCost.Minutes() >= minCost.Minutes() ||
			graph.Nodes[dest.DestinationNodeId].VisitStatus == entity.VisitStatusVisited {
			continue
		}
		calculateOptimalRoute(dest.DestinationNodeId, newCost, graph, minCost, bestRoute, curRoute)
	}
	node.VisitStatus = entity.VisitStatusUnVisited
	graph.Nodes[src] = node
	*curRoute = (*curRoute)[:len(*curRoute)-1]
}
