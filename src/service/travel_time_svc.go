package service

import (
	"errors"
	"fmt"
	"log"
	"microservice/src/config"
	"microservice/src/entity"
	"microservice/src/helper"
	"time"
)

type (
	TravelTimeServiceImpl struct {
		cfg config.Config
	}
)

func NewTravelTimeSvc(cfg config.Config) TravelTimeSvc {
	return &TravelTimeServiceImpl{
		cfg: cfg,
	}
}

func (s TravelTimeServiceImpl) CalculateTravelTime(agent entity.DeliveryAgent,
	distance entity.Distance) (time.Duration, error) {
	var duration time.Duration

	if agent.AverageSpeed.Value <= 0 {
		err := errors.New("invalid agent speed")

		log.Printf("[ERROR] invalid agent speed: %s", err.Error())

		return duration, err
	}

	travelTime := distance.Value / agent.AverageSpeed.Value
	travel := helper.RoundTo(travelTime, s.cfg.CalculationConfig.DistancePrecisionDecimals)
	durationString := fmt.Sprintf("%f%s", travel, s.cfg.CalculationConfig.DistanceUnit)
	duration, err := time.ParseDuration(durationString)

	if err != nil {
		log.Printf("[ERROR] error parsing travel duration: %s", err.Error())
		return duration, err
	}

	return duration, nil
}
