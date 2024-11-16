package app

import (
	"fmt"
	"os"
)

func AppHome() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appHome := fmt.Sprintf("%s/GoSkyScheduler", configDir)

	_, homeExistsErr := os.Stat(appHome)
	if homeExistsErr != nil {
		mkdirErr := os.MkdirAll(appHome, os.ModeDir)

		if mkdirErr != nil {
			return "", mkdirErr
		}

	}

	return appHome, nil
}
