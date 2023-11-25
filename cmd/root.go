package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var envName string
var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "envsync",
	Short: "Sync your .env files securely.",
	Long:  `Sync your .env files securely across multiple parties using private/public key!`,
}

func Execute() {
	cobra.OnInitialize(initConfig)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.envsync.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home + "/.envsync")
		viper.SetConfigName("envsync")
	}

	viper.AutomaticEnv()
}
