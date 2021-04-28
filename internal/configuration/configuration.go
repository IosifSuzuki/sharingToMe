package configuration

import (
	"IosifSuzuki/sharingToMe/internal/models"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

const (
	nameOfConfigFile = "info"
	typeOfConfigurationFile = "yml"
)

var Configuration = setupConfiguration()

func setupConfiguration() *models.ConfigurationFile {
	baseDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	viper.AddConfigPath(filepath.Join(baseDir, "src", "config"))
	viper.SetConfigName(nameOfConfigFile)
	viper.SetConfigType(typeOfConfigurationFile)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	var configuration models.ConfigurationFile
	if err := viper.Unmarshal(&configuration); err != nil {
		panic(err)
	}

	return &configuration
}
