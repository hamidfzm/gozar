package cmd

import (
	"fmt"

	"gozar/config"

	"github.com/spf13/cobra"
)

var Version string

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: fmt.Sprintf("Print the version number of %s server", config.Name),
	Run: func(cmd *cobra.Command, args []string) {
		if Version == "" {
			Version = "v0.0.0"
		}

		fmt.Printf("%s server %s\n", config.Name, Version)
	},
}
