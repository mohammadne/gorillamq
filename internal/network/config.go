package network

type Config struct {
	// default pool of connections

	// https://blog.pinterjann.is/ed25519-certificates.html
	TLS struct {
		Enabled     bool   `koanf:"enabled"`
		Certificate string `koanf:"certificate"`
		PrivateKey  string `koanf:"private_key"`
	} `koanf:"tls"`
}
