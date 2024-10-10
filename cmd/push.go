/*
Copyright © 2024 Jesse Maitland jesse@pytoolbelt.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/pytoolbelt/ime/pkg/config"
	"github.com/pytoolbelt/ime/pkg/environment"
	"github.com/pytoolbelt/ime/pkg/paramstore"
	"github.com/spf13/cobra"
)

var modeFlag string
var envFlag string
var projFlag string
var overwriteFlag bool

func AddParameters(ps *paramstore.ParamStore, ef *environment.EnvFile) error {
	fmt.Print("adding parameters to parameter store")
	return nil
}

func DeleteParameters(ps *paramstore.ParamStore, ef *environment.EnvFile) error {
	fmt.Print("deleting parameters from parameter store")
	return nil
}

func MergeParameters(ps *paramstore.ParamStore, ef *environment.EnvFile) error {
	fmt.Print("merging parameters with parameter store")
	return nil
}

func IsValidMode(mode string) bool {
	switch mode {
	case "add", "delete", "merge":
		return true
	default:
		return false
	}
}

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push an environment to AWS Parameter Store",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

		if !IsValidMode(modeFlag) {
			fmt.Printf("Invalid mode: %s \n", modeFlag)
			os.Exit(1)
		}

		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading configuration: %s \n", err)
			os.Exit(1)
		}

		envConf, err := cfg.GetEnvironment(projFlag, envFlag)
		if err != nil {
			fmt.Printf("Error getting environment from config ime.yaml: %s \n", err)
			os.Exit(1)
		}

		ef := environment.NewEnvFileFromPath(envConf.GetResolvedLocalPath())
		if err := ef.LoadEnvFile(); err != nil {
			fmt.Printf("Error loading environment file: %s \n", err)
			os.Exit(1)
		}

		psPath, err := cfg.FormatParameterStorePath(projFlag, envFlag)
		if err != nil {
			fmt.Printf("Error formatting parameter store path: %s \n", err)
			os.Exit(1)
		}

		ps, err := paramstore.NewParamStore(psPath)
		if err != nil {
			fmt.Printf("Error creating ParamStore: %s \n", err)
			os.Exit(1)
		}

		switch modeFlag {

		case "add":
			fmt.Printf("Adding parameters to %s \n", psPath)
			if err := AddParameters(ps, ef); err != nil {
				fmt.Printf("Error adding parameters: %s \n", err)
				os.Exit(1)
			}
		case "delete":
			fmt.Printf("Deleting parameters from %s \n", psPath)
			if err := DeleteParameters(ps, ef); err != nil {
				fmt.Printf("Error deleting parameters: %s \n", err)
				os.Exit(1)
			}
		case "merge":
			fmt.Printf("Merging parameters with %s \n", psPath)
			if err := MergeParameters(ps, ef); err != nil {
				fmt.Printf("Error merging parameters: %s \n", err)
				os.Exit(1)
			}
		default:
			fmt.Printf("Invalid mode: %s \n", modeFlag)
			os.Exit(1)
		}

		// eParams, err := ps.GetParameters()
		// if err != nil {
		// 	fmt.Printf("Error getting parameters: %s \n", err)
		// 	os.Exit(1)
		// }

	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().StringVar(&projFlag, "project", "", "The project to push")
	pushCmd.Flags().StringVar(&envFlag, "env", "", "The environment to push")
	pushCmd.Flags().StringVar(&modeFlag, "mode", "add", "Mode of operation: add, delete, or merge")

	pushCmd.Flags().BoolVar(&overwriteFlag, "overwrite", false, "Overwrite existing parameters in Parameter Store")

	// Mark the required flags
	pushCmd.MarkFlagRequired("project")
	pushCmd.MarkFlagRequired("env")
}
