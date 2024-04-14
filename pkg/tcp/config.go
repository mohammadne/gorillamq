package tcp

type Config struct {
	SecurePort   int `koanf:"secure_port"`
	InsecurePort int `koanf:"insecure_port"`

	// default pool of connections

	// https://blog.pinterjann.is/ed25519-certificates.html
	TLS struct {
		Certificate string `koanf:"certificate"`
		PrivateKey  string `koanf:"private_key"`
	} `koanf:"tls"`
}
