package commons

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CheckConfig(cmd *cobra.Command, args []string) {
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("No config file found. Please configure the config by running init or provide the config with --config")
		os.Exit(1)
	}

	requiredKeys := []string{
		"aws.access_key_id",
		"aws.region",
		"aws.s3_bucket",
		"aws.secret_access_key",
		"envsync.private_key",
		"envsync.public_key",
	}

	for _, key := range requiredKeys {
		if !viper.IsSet(key) {
			fmt.Printf("Required configuration key missing: %s\n", key)
			os.Exit(1)
		}
	}

	fmt.Println("Using config file:", viper.ConfigFileUsed())
}
