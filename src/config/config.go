package config

type (
	CalculationConfig struct {
		DistancePrecisionDecimals int     `yaml:"distance_precision_decimals"`
		DistanceUnit              string  `yaml:"distance_unit"`
		RadiusKm                  float64 `yaml:"radius_km"`
	}

	Config struct {
		CalculationConfig CalculationConfig `yaml:"calculation"`
	}
)
