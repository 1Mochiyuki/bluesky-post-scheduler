package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

type Config struct {
	env map[string]string
}

func appHome() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Error("failed to get user config dir", err)
		panic("failed to get user config dir")
	}
	appHome := fmt.Sprintf("%s/GoSkyScheduler", configDir)
	log.Debug("app home", appHome)

	_, homeExistsErr := os.Stat(appHome)
	if homeExistsErr != nil {
		mkdirErr := os.MkdirAll(appHome, os.ModeDir)

		if mkdirErr != nil {
			return ""
		}

		log.Debug("home created at", appHome)
	}

	return appHome
}

func initConfig() error {
	viper.SetConfigName("vars")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("GSKY")
	viper.BindEnv("app_pass")
	viper.SetDefault("gsky_app_pass", "")
	log.Info("here")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error("err", err)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Error("Config file not found", err)
			safeWriteErr := viper.SafeWriteConfig()
			if safeWriteErr != nil {
				log.Error("safe write err", safeWriteErr)
				return safeWriteErr
			}
			log.Info("config file created")
		}
	}
	viper.WatchConfig()

	return nil
}
