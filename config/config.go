// Viper library is used here to extract environment variables.
// This is done because it is wildly popular and provides convenient, extensible functionality
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Server struct {
	Host   string `mapstructure:"HOST"`
	Port   int    `mapstructure:"PORT"`
	ApiKey string `mapstructure:"API_KEY"`
}

type Database struct {
	Host     string `mapstructure:"HOST"`
	Port     int    `mapstructure:"PORT"`
	User     string `mapstructure:"USER"`
	Password string `mapstructure:"PASSWORD"`
	DB       string `mapstructure:"DB"`
}

type Calculator struct {
	TimeAdditionMs        int `mapstructure:"TIME_ADDITION_MS"`
	TimeSubtractionMs     int `mapstructure:"TIME_SUBTRACTION_MS"`
	TimeMultiplicationsMs int `mapstructure:"TIME_MULTIPLICATIONS_MS"`
	TimeDivisionsMs       int `mapstructure:"TIME_DIVISIONS_MS"`
}

type Config struct {
	Server           Server     `mapstructure:"SERVER"`
	Database         Database   `mapstructure:"DATABASE"`
	CalculatorConfig Calculator `mapstructure:"CALC"`
}

func validatePort(port int) bool {
	return port >= 0 && port <= 65535
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(fmt.Sprintf("%s/.env", path))

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	if !validatePort(config.Server.Port) {
		err = fmt.Errorf("invalid port: %d", config.Server.Port)
		return
	}

	return
}
