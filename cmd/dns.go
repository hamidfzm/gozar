package cmd

import (
	"strings"

	"gozar/config"
	"gozar/dns"
	"gozar/dns/handlers"

	"github.com/spf13/cobra"

	"github.com/getsentry/sentry-go"
)

func init() {
	rootCmd.AddCommand(dnsCmd)
}

func configureDNSServer() (*dns.Server, error) {
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

	dnsServer := dns.New(c)
	err := handlers.Configure(dnsServer)
	if err != nil {
		return nil, err
	}

	return dnsServer, nil
}

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "Run dns server",
	RunE: func(cmd *cobra.Command, args []string) error {
		dnsServer, err := configureDNSServer()
		if err != nil {
			return err
		}

		return dnsServer.Start()
	},
}
