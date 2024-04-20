package service

import (
	"microservice/src/config"
	"microservice/src/entity"
	"microservice/src/helper"
)

type (
	HaversineDistanceSvc struct {
		cfg config.Config
	}
)

func NewHaversineDistanceSvc(cfg config.Config) DistanceSvc {
	return &HaversineDistanceSvc{
		cfg: cfg,
	}
}

func (s HaversineDistanceSvc) CalculateDistance(origin entity.GeoCoordinates,
	destination entity.GeoCoordinates) entity.Distance {
	centralAngleRadians := helper.CalculateCentralAngleRadians(origin, destination)

	return entity.Distance{
		Value: centralAngleRadians * s.cfg.CalculationConfig.RadiusKm,
		Unit:  entity.DistanceUnitKm,
	}
}
