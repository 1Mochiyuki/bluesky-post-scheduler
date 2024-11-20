package app

import (
	"os"
	"path/filepath"
)

func AppHome() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	// appHome := fmt.Sprintf("%s/GoSkyScheduler.", configDir)

	appHome := filepath.Join(configDir, "GoSkyScheduler")

	_, homeExistsErr := os.Stat(appHome)
	if homeExistsErr != nil {
		mkdirErr := os.MkdirAll(appHome, os.ModeDir)

		if mkdirErr != nil {
			return "", mkdirErr
		}

	}

	return appHome, nil
}
