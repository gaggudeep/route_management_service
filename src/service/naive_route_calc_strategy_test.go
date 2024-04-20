package service

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"microservice/src/entity"
	"testing"
	"time"
)

type travelSvcMock struct{}
type distanceSvcMock struct{}

func (d distanceSvcMock) CalculateDistance(origin entity.GeoCoordinates,
	dest entity.GeoCoordinates) entity.Distance {
	return entity.Distance{
		Value: math.Abs(origin.Longitude - dest.Longitude),
		Unit:  entity.DistanceUnitKm,
	}
}

func (t travelSvcMock) CalculateTravelTime(agent entity.DeliveryAgent,
	distance entity.Distance) (time.Duration, error) {
	return time.ParseDuration(fmt.Sprintf("%fh", distance.Value/agent.AverageSpeed.Value))
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

func TestNaiveRouteCalcStrategy_CalculateRoute(t *testing.T) {
	orderSvcMock := OrderGraphSvcImpl{
		travelTimeSvc: travelSvcMock{},
		distanceSvc:   distanceSvcMock{},
	}
	svc := NewNaiveRouteCalcStrategy(orderSvcMock)

	type args struct {
		agent  entity.DeliveryAgent
		orders []entity.Order
	}
	tests := []struct {
		name      string
		args      args
		wantError bool
		expected  entity.Route
		before    func(*testing.T)
	}{
		{
			name: "returns optimal route",
			args: args{
				agent: agent,
				orders: []entity.Order{
					{
						Restaurant: res1,
						Consumer:   con1,
					},
					{
						Restaurant: res2,
						Consumer:   con2,
					},
				},
			},
			expected: entity.Route{
				RouteNodes: []entity.RouteNode{
					{
						OriginalNode: agent,
					},
					{
						MinutesToReachFromStart: 30,
						OriginalNode:            res2,
					},
					{
						MinutesToReachFromStart: 51.9,
						OriginalNode:            con2,
					},
					{
						MinutesToReachFromStart: 73.8,
						OriginalNode:            res1,
					},
					{
						MinutesToReachFromStart: 104.7,
						OriginalNode:            con1,
					},
				},
			},
		},
		{
			name: "returns empty route",
			before: func(t *testing.T) {
				t.Skip()
			},
		},
		{
			name: "given multiple orders in same restaurant, returns optimal route",
			before: func(t *testing.T) {
				t.Skip()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(t)
			}
			actual, err := svc.CalculateRoute(tt.args.agent, tt.args.orders)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			s, err := json.Marshal(actual)
			assert.NoError(t, err)

			t.Logf("json: %s", string(s))
			assert.Equal(t, tt.expected, actual)
		})
	}
}
