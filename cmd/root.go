package cmd

import (
	"fmt"
	"os"

	"gozar/config"

	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(config.InitViperConfig)
}

var rootCmd = &cobra.Command{
	Use:   config.Name,
	Short: fmt.Sprintf("%s server", config.Name),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
