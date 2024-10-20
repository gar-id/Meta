/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"MetaHandler/agent/windows/setup"

	"MetaHandler/core/caches"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Version: caches.MetaHandlerVersion,
	Use:     "install",
	Short:   "Install Meta-Handler Agent",
	Long: `Install Meta-Handler Agent to server Windows only. 
You can not install this handler Into non-Windows
OS. Any unknown behavior will be ignored`,
	Run: func(cmd *cobra.Command, args []string) {
		setup.Install()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
