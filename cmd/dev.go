package cmd

import (
	"github.com/spf13/cobra"

	"os"
	"os/signal"
)

func init() {
	rootCmd.AddCommand(devCmd)
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Run everything for dev",
	RunE: func(cmd *cobra.Command, args []string) error {
		dnsServer, err := configureDNSServer()
		if err != nil {
			return err
		}

		proxyServer, err := configureProxyServer()
		if err != nil {
			return err
		}

		go func() {
			err = dnsServer.Start()
			if err != nil {
				panic(err)
			}
		}()

		go func() {
			err = proxyServer.Start()
			if err != nil {
				panic(err)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit
		if err := dnsServer.Shutdown(); err != nil {
			return err
		}

		return nil
	},
}
