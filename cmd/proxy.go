package cmd

import (
	"strings"

	"gozar/config"
	"gozar/proxy"

	"github.com/spf13/cobra"

	"github.com/getsentry/sentry-go"
)

func init() {
	rootCmd.AddCommand(proxyCmd)
}

func configureProxyServer() (*proxy.Server, error) {
	c := config.NewViperConfig()

	if c.Sentry.Enable {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:   c.Sentry.DSN,
			Debug: strings.EqualFold(c.Log.Level, "DEBUG"),
		})
		if err != nil {
			return nil, err
		}
	}

	proxyServer := proxy.New(c)

	return proxyServer, nil
}

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Run proxy server",
	RunE: func(cmd *cobra.Command, args []string) error {
		proxyServer, err := configureProxyServer()
		if err != nil {
			return err
		}

		return proxyServer.Start()
	},
}
