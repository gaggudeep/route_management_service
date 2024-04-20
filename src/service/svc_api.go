package service

import (
	"microservice/src/entity"
	"time"
)

type (
	RouteSvc interface {
		// CalculateRoute NOTE: restaurants and consumers must be linked by having same slice index.
		FindRoute(entity.DeliveryAgent, []entity.Order) (entity.Route, error)
	}

	DistanceSvc interface {
		// CalculateDistance NOTE: for now only supports kkm distance
		CalculateDistance(entity.GeoCoordinates, entity.GeoCoordinates) entity.Distance
	}

	TravelTimeSvc interface {
		// CalculateTravelTime NOTE: for now only support km distance and hour time unit for calculation
		CalculateTravelTime(entity.DeliveryAgent, entity.Distance) (time.Duration, error)
	}

	RouteCalculationStrategy interface {
		CalculateRoute(entity.DeliveryAgent, []entity.Order) (entity.Route, error)
	}

	OrderGraphSvc interface {
		ConstructGraph(entity.DeliveryAgent, []entity.Order) (entity.Graph[time.Duration], error)
	}
)
