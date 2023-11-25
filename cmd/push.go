package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/setkyar/envsync/commons"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:    "push",
	Short:  "Push your .env file",
	Long:   `This command will push your .env file from current directory.`,
	PreRun: commons.CheckConfig,
	Run: func(cmd *cobra.Command, args []string) {
		envPath, _ := cmd.Flags().GetString("env")
		if envPath == "" {
			defaultEnvPath := ".env"
			if _, err := os.Stat(defaultEnvPath); os.IsNotExist(err) {
				fmt.Println("No .env file found in the current directory. Please specify a path using --env flag.")
				return
			}
			envPath = defaultEnvPath
		}

		absPath, err := filepath.Abs(envPath)
		if err != nil {
			fmt.Printf("Error resolving absolute path: %s\n", err)
			return
		}

		if err := commons.UploadToS3(envName, absPath); err != nil {
			fmt.Println("Error uploading to S3:", err)
			return
		}
	},
}

func init() {
	pushCmd.Flags().StringVarP(&envName, "name", "n", "", "Environment name (required)")
	pushCmd.MarkFlagRequired("name")

	pushCmd.Flags().StringP("env", "e", "", "Optional custom path for the .env file")

	rootCmd.AddCommand(pushCmd)
}
