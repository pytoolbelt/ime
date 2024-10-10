// /*
// Copyright © 2024 Jesse Maitland jesse@pytoolbelt.com
// */
package cmd

// import (
// 	"fmt"
// 	"os"

// 	"github.com/pytoolbelt/ime/pkg/config"
// 	"github.com/pytoolbelt/ime/pkg/environment"
// 	"github.com/pytoolbelt/ime/pkg/paramstore"
// 	"github.com/spf13/cobra"
// )

// // fetchCmd represents the fetch command
// var fetchCmd = &cobra.Command{
// 	Use:   "fetch",
// 	Short: "fetch an environment from the AWS Parameter Store",
// 	Long:  "fetches an environment from the AWS Parameter Store and saves it in the .env file.",
// 	Run: func(cmd *cobra.Command, args []string) {

// 		cfg, err := config.LoadConfig()
// 		if err != nil {
// 			fmt.Printf("Error loading configuration: %s \n", err)
// 			os.Exit(1)
// 		}

// 		projectName, err := cmd.Flags().GetString("project")
// 		if err != nil {
// 			fmt.Printf("Error getting project: %s \n", err)
// 			os.Exit(1)
// 		}

// 		environmentName, err := cmd.Flags().GetString("env")
// 		if err != nil {
// 			fmt.Printf("Error getting environment: %s \n", err)
// 			os.Exit(1)
// 		}

// 		path, err := cfg.FormatParameterStorePath(projectName, environmentName)
// 		if err != nil {
// 			fmt.Printf("Error formatting parameter store path: %s \n", err)
// 			os.Exit(1)
// 		}

// 		env, err := cfg.GetEnvironment(projectName, environmentName)
// 		if err != nil {
// 			fmt.Printf("Error getting environment: %s \n", err)
// 			os.Exit(1)
// 		}

// 		fmt.Printf("Fetching %s from %s \n", environmentName, path)

// 		result, err := paramstore.GetParametersByPath(path)

// 		if err != nil {
// 			fmt.Printf("Error fetching parameters: %s \n", err)
// 			os.Exit(1)
// 		}

// 		params, err := paramstore.ParseParameterKeyValuePairs(*result)

// 		if err != nil {
// 			fmt.Printf("Error parsing parameters: %s \n", err)
// 			os.Exit(1)
// 		}

// 		ef := environment.NewEnvFileFromPath(env.GetResolvedLocalPath())
// 		if err := ef.LoadEnvFile(); err != nil {
// 			fmt.Printf("Error loading environment file: %s \n", err)
// 			os.Exit(1)
// 		}

// 		ef.ConvertEnvVarsToMap(params)

// 		if err := ef.WriteEnvFile(); err != nil {
// 			fmt.Printf("Error writing environment file: %s \n", err)
// 			os.Exit(1)
// 		}
// 	},
// }

// func init() {
// 	rootCmd.AddCommand(fetchCmd)
// 	fetchCmd.Flags().String("project", "", "The project to fetch")
// 	fetchCmd.Flags().String("env", "", "The project environment to fetch")

// 	// Mark the required flags
// 	fetchCmd.MarkFlagRequired("project")
// 	fetchCmd.MarkFlagRequired("env")

// 	// init viper config
// 	config.InitializeConfig()
// }
