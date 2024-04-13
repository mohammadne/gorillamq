package auth

type Config struct {
	Username string `koanf:"username"`
	Password string `koanf:"password"`
}
