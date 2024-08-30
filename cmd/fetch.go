/*
Copyright © 2024 Jesse Maitland jesse@pytoolbelt.com
*/
package cmd

import (
	"github.com/pytoolbelt/ime/pkg/terminal"
	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch an environment from the AWS Parameter Store",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		terminal.StartSubshell()
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fetchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
