package utils

import "github.com/spf13/viper"

type Config struct {
	DBdriver  string `mapstructure:"DB_DRIVER"`
	DB_source string `mapstructure:"DB_SOURCE"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		return nil, err
	}

	var config Config

	err = viper.Unmarshal(&config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}