package config

type DNS struct {
	Host string
	Port uint16

	Domains []string
	Proxies []string
}
