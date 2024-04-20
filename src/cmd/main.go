package main

import (
	"encoding/json"
	"log"
	"microservice/src/config"
	"microservice/src/entity"
	"microservice/src/service"
)

func main() {
	cfg := config.Config{
		CalculationConfig: config.CalculationConfig{
			DistancePrecisionDecimals: 2,
			DurationUnit:              "m",
			RadiusKm:                  6378,
		},
	}
	travelTimeSvc := service.NewTravelTimeSvc(cfg)
	distSvc := service.NewHaversineDistanceSvc(cfg)
	orderGraphSvc := service.NewOrderGraphSvvImpl(travelTimeSvc, distSvc)
	routeCalcStaretegy := service.NewNaiveRouteCalcStrategy(orderGraphSvc)
	routeSvc := service.NewRouteSvc(routeCalcStaretegy)

	route, err := routeSvc.FindRoute(agent, []entity.Order{
		{
			Restaurant: res1,
			Consumer:   con1,
		},
		{
			Restaurant: res2,
			Consumer:   con2,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	jsonRes, err := json.MarshalIndent(route, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(jsonRes))
}

var (
	agent = entity.DeliveryAgent{
		AverageSpeed: entity.Speed{
			Value:        20,
			DistanceUnit: entity.DistanceUnitKm,
			TimeUnit:     entity.TimeUnitHr,
		},
		Location: entity.Location{
			Type: entity.NodeTypeAgent,
		},
	}
	res1 = entity.Restaurant{
		Location: entity.Location{
			GeoCoordinates: entity.GeoCoordinates{
				Latitude:  10,
				Longitude: 10,
			},
			Id:   1,
			Type: entity.NodeTypeRestaurant,
		},
		MeanPreparationTime: 20,
	}
	con1 = entity.Consumer{
		Location: entity.Location{
			GeoCoordinates: entity.GeoCoordinates{
				Latitude:  20.5,
				Longitude: 20.3,
			},
			Id:   1,
			Type: entity.NodeTypeConsumer,
		},
	}
	res2 = entity.Restaurant{
		Location: entity.Location{
			GeoCoordinates: entity.GeoCoordinates{
				Latitude:  10,
				Longitude: 10,
			},
			Id:   2,
			Type: entity.NodeTypeRestaurant,
		},
		MeanPreparationTime: 40,
	}
	con2 = entity.Consumer{
		Location: entity.Location{
			GeoCoordinates: entity.GeoCoordinates{
				Latitude:  12.5,
				Longitude: 17.3,
			},
			Id:   2,
			Type: entity.NodeTypeConsumer,
		},
	}
)
