package config

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"os"
)

// initConfig reads in config file and ENV variables if set.
func InitConfig(cfg string) func() {
	return func() {

		if cfg != "" {
			// Use config file from the flag.
			viper.SetConfigFile(cfg)
		} else {
			// Find home directory.
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// Search config in home directory with name ".goscript" (without extension).
			viper.AddConfigPath(home)
			viper.SetConfigName(".chronic")
		}

		viper.AutomaticEnv() // read in environment variables that match

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}
}
