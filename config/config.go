package config

import (
	"github.com/1Mochiyuki/gosky/config/logger"
	"github.com/spf13/viper"
)

type Config struct {
	env map[string]string
}

var l = logger.Get()

func InitConfig() error {
	viper.SetConfigName("vars")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("GSKY")
	viper.BindEnv("app_pass")
	viper.SetDefault("gsky_app_pass", "")
	err := viper.ReadInConfig()
	if err != nil {
		l.Error().Err(err).Msg("received error")
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			safeWriteErr := viper.SafeWriteConfig()
			if safeWriteErr != nil {
				l.Error().Err(safeWriteErr).Msg("safe write err")
				return safeWriteErr
			}
			l.Info().Msg("config file created")
		}
	}
	viper.WatchConfig()

	return nil
}
