package service

import (
	"log"
	"microservice/src/entity"
)

type (
	RouteSvcImpl struct {
		routeCalculationStrategy RouteCalculationStrategy
	}
)

func NewRouteSvc(routeCalculationStrategy RouteCalculationStrategy) RouteSvc {
	return &RouteSvcImpl{
		routeCalculationStrategy: routeCalculationStrategy,
	}
}

func (r RouteSvcImpl) FindRoute(agent entity.DeliveryAgent, orders []entity.Order) (entity.Route, error) {
	route, err := r.routeCalculationStrategy.CalculateRoute(agent, orders)
	if err != nil {
		log.Printf("[ERROR] error calculating route : %s", err.Error())
		return route, err
	}

	return route, err
}
