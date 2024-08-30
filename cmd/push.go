/*
Copyright © 2024 Jesse Maitland jesse@pytoolbelt.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/pytoolbelt/ime/pkg/config"
	"github.com/pytoolbelt/ime/pkg/paramstore"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push an environment to AWS Parameter Store",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()

		if err != nil {
			fmt.Printf("Error loading configuration: %s \n", err)
			os.Exit(1)
		}

		projectName, err := cmd.Flags().GetString("project")
		if err != nil {
			fmt.Printf("Error getting project: %s \n", err)
			os.Exit(1)
		}

		environmentName, err := cmd.Flags().GetString("env")
		if err != nil {
			fmt.Printf("Error getting environment: %s \n", err)
			os.Exit(1)
		}

		overwrite, err := cmd.Flags().GetBool("overwrite")
		if err != nil {
			fmt.Printf("Error getting overwrite flag: %s \n", err)
			os.Exit(1)
		}

		env, err := cfg.GetEnvironment(projectName, environmentName)
		if err != nil {
			fmt.Printf("Error getting environment: %s \n", err)
			os.Exit(1)
		}

		envFile, err := paramstore.LoadEnvFile(env.GetResolvedLocalPath())
		if err != nil {
			fmt.Printf("Error loading environment file: %s \n", err)
			os.Exit(1)
		}

		paramStorePath, err := cfg.FormatParameterStorePath(projectName, environmentName)
		if err != nil {
			fmt.Printf("Error formatting parameter store path: %s \n", err)
			os.Exit(1)
		}

		err = paramstore.PutParametersFromEnvFile(envFile, paramStorePath, overwrite)
		if err != nil {
			fmt.Printf("Error putting parameters: %s \n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().String("project", "", "The project to push")
	pushCmd.Flags().String("env", "", "The environment to push")
	pushCmd.Flags().BoolP("overwrite", "o", false, "Overwrite existing parameters in Parameter Store")

	// Mark the required flags
	pushCmd.MarkFlagRequired("project")
	pushCmd.MarkFlagRequired("env")
}
