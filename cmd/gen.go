package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	genCmd.AddCommand(confCmd)
	rootCmd.AddCommand(genCmd)
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate samples",
}

var confCmd = &cobra.Command{
	Use:   "config",
	Short: "Generate basic configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Generating config.yml")
		return viper.WriteConfigAs("config.yml")
	},
}
