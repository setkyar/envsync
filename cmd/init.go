package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/setkyar/envsync/commons"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize your envsync",
	Long:  `This command will setup sqlite/ public / private key and AWS S3.`,
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()

		if err != nil {
			fmt.Println("Failed to get home directory")
			return
		}

		keyDir := filepath.Join(homeDir, ".envsync")
		privateKeyPath := filepath.Join(keyDir, "private_key.pem")
		publicKeyPath := filepath.Join(keyDir, "public_key.pem")
		configPath := filepath.Join(keyDir, "envsync.yaml")

		if _, err := os.Stat(configPath); err == nil {
			fmt.Println("You have already configured. Please update your config manually at " + configPath)
			return
		}

		fmt.Println("Initializing envsync...")

		awsConfig := make(map[string]string)
		inputFields := []string{"AWS Region", "AWS Access Key ID", "AWS Secret Access Key", "S3 Bucket Name"}

		reader := bufio.NewReader(os.Stdin)
		for _, field := range inputFields {
			fmt.Printf("%s: ", field)
			input, _ := reader.ReadString('\n')

			// Trim space and newline character
			input = strings.TrimSpace(input)
			awsConfig[field] = input
		}

		viper.Set("aws.region", awsConfig["AWS Region"])
		viper.Set("aws.access_key_id", awsConfig["AWS Access Key ID"])
		viper.Set("aws.secret_access_key", awsConfig["AWS Secret Access Key"])
		viper.Set("aws.s3_bucket", awsConfig["S3 Bucket Name"])
		viper.Set("envsync.private_key", privateKeyPath)
		viper.Set("envsync.public_key", publicKeyPath)

		// Save the configuration to a file
		viper.SetConfigType("yaml")
		viper.SetConfigFile(configPath)
		if err := viper.WriteConfig(); err != nil {
			fmt.Printf("Failed to write configuration: %s\n", err)
			return
		}

		fmt.Println("Configuration saved successfully at " + configPath)

		if _, err := os.Stat(privateKeyPath); err == nil {
			fmt.Println("Private key already exists.")
			return
		}

		if _, err := os.Stat(publicKeyPath); err == nil {
			fmt.Println("Public key already exists.")
			return
		}

		fmt.Println("Generating RSA public/private keys...")
		if err := commons.GenerateRSAKeys(); err != nil {
			fmt.Printf("Failed to generate RSA keys: %s\n", err)
			return
		}
		fmt.Println("RSA keys generated successfully.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
