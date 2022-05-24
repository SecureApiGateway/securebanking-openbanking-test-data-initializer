package common

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strconv"
)

func LoadConfigurationByEnv(environment string) (err error) {
	// name of config file (without extension)
	viper.SetConfigName("viper-" + environment + "-configuration")
	// REQUIRED if the config file does not have the extension in the name
	viper.SetConfigType("yaml")
	// call multiple times to add many search paths
	viper.AddConfigPath("config/viper")
	// optionally look for config in the working directory
	viper.AddConfigPath(".")
	// Viper will check for an environment variable any time a viper.Get request is made
	viper.AutomaticEnv()
	// Find and read the config file
	err = viper.ReadInConfig()
	// Handle errors reading the config file
	if err != nil {
		//panic(fmt.Errorf("Fatal error config file: %w \n", err))
		return
	}
	err = viper.Unmarshal(&Config)
	fmt.Println("Verbose / Debug level [", strconv.FormatBool(Config.Environment.Verbose), "]")
	return
}

func ConfigureLogger() (*zap.Logger, error) {
	verbose := viper.GetBool("ENVIRONMENT.VERBOSE")

	if verbose {
		configZap := zap.NewProductionConfig()
		configZap.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		// disable sampling to ensure we get all log messages
		configZap.Sampling = nil
		return configZap.Build(zap.AddCaller())
	}
	return zap.NewProduction(zap.AddCaller())
}
