package settings

// ClientConfig represents only client configuration
type ClientConfig struct {
	ServerAddress string
	UseTLS        bool
}

func GetClientConfig() ClientConfig {
	var config ClientConfig

	flags := Parse()

	config.ServerAddress = (*flags).Endpoint
	config.UseTLS = false

	return config
}
