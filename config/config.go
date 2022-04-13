package config

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const Name = "gozar"

type Config struct {
	Log    *Log
	DNS    *DNS
	Proxy  *Proxy
	Sentry *Sentry
}

func InitViperConfig() {
	viper.SetDefault("log", map[string]interface{}{
		"level": "debug",
	})
	viper.SetDefault("dns", map[string]interface{}{
		"host": "0.0.0.0",
		"port": 53,

		"domains": []string{
			"google.com",
		},

		"proxies": []string{
			"94.130.124.121",
		},
	})
	viper.SetDefault("proxy", map[string]interface{}{
		"host":  "0.0.0.0",
		"ports": []uint16{80, 443},
	})
	viper.SetDefault("sentry", map[string]interface{}{
		"dsn":    "",
		"enable": false,
	})

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix(Name)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	switch err.(type) {
	case viper.ConfigFileNotFoundError, nil:
	default:
		panic(err)
	}
}

func NewViperConfig() *Config {
	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	switch strings.ToUpper(cfg.Log.Level) {
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
	default:
		logrus.SetLevel(logrus.WarnLevel)
	}

	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	return &cfg
}
