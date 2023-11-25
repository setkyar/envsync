package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/setkyar/envsync/commons"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:    "pull",
	Short:  "Pull your environment file",
	Long:   `This command will pull your .env file to current directory`,
	PreRun: commons.CheckConfig,
	Run: func(cmd *cobra.Command, args []string) {
		envPath, _ := cmd.Flags().GetString("env")

		if envPath == "" {
			envPath = ".env"

		}
		absPath, err := filepath.Abs(envPath)

		if err != nil {
			fmt.Printf("Error resolving absolute path: %s\n", err)
			return
		}

		if err := commons.DownloadFromS3(envName, absPath); err != nil {
			fmt.Println("Error downloading from S3:", err)
			return
		}
	},
}

func init() {
	pullCmd.Flags().StringVarP(&envName, "name", "n", "", "Environment name (required)")
	pullCmd.MarkFlagRequired("name")

	pullCmd.Flags().StringP("env", "e", "", "Optional custom path for the .env file")
	rootCmd.AddCommand(pullCmd)
}
