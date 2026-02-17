package config

import "github.com/spf13/viper"

type Config struct {
	HTTPPort string

	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	DBSSLMode string
}

func Load() (Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	cfg := Config{
		HTTPPort: viper.GetString("HTTP_PORT"),

		DBHost:    viper.GetString("DB_HOST"),
		DBPort:    viper.GetString("DB_PORT"),
		DBUser:    viper.GetString("DB_USER"),
		DBPass:    viper.GetString("DB_PASSWORD"),
		DBName:    viper.GetString("DB_NAME"),
		DBSSLMode: viper.GetString("DB_SSLMODE"),
	}

	return cfg, nil
}
