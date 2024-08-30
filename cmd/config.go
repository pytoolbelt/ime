/*
Copyright © 2024 Jesse Maitland jesse@pytoolbelt.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/pytoolbelt/ime/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func configEntrypoint(cmd *cobra.Command, args []string) {
	config.InitializeConfig()

	switch {
	case cmd.Flags().Changed("show"):
		printConfig()

	case cmd.Flags().Changed("path"):
		printPath()

	default:
		fmt.Println("No action specified for the config command")
		fmt.Printf("Must pass one of the following flags: \n%s", cmd.Flags().FlagUsages())
	}
}

func printConfig() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading configuration: %s", err)
		os.Exit(1)
	}

	configFile := viper.ConfigFileUsed()
	fmt.Printf("Config file used: %s\n", configFile)

	cfg.PrintTable()
	os.Exit(0)
}

func printPath() {
	configFile := viper.ConfigFileUsed()
	fmt.Printf("Config file used: %s\n", configFile)
	os.Exit(0)
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Interact with the ime configuration file",
	Long:  "",
	Run:   configEntrypoint,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().BoolP("show", "s", false, "List the current configuration")
	configCmd.Flags().BoolP("path", "p", false, "Show the path to the configuration file used")
}
