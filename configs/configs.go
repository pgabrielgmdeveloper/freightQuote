package configs

import "github.com/spf13/viper"

type conf struct {
	DBDriver         string `mapstructure:"DB_DRIVER"`
	DBHost           string `mapstructure:"DB_HOST"`
	DBPort           string `mapstructure:"DB_PORT"`
	DBUser           string `mapstructure:"DB_USER"`
	DBPassword       string `mapstructure:"DB_PASSWORD"`
	DBName           string `mapstructure:"DB_NAME"`
	TokenAPI         string `mapstructure:"TOKEN_API"`
	PlatformCode     string `mapstructure:"PLATFORM_CODE"`
	RegisteredNumber string `mapstructure:"REGISTERED_NUMBER"`
}

func LoadConfig() (*conf, error) {
	var cfg *conf
	viper.AutomaticEnv()
	viper.BindEnv("DB_DRIVER")
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("TOKEN_API")
	viper.BindEnv("PLATFORM_CODE")
	viper.BindEnv("REGISTERED_NUMBER")
	err := viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
