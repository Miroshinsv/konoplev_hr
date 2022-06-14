package config

type JWT struct {
	Secret       string `env:"JWT_SECRET" envDefault:"default"`
	TokenVersion uint   `env:"JWT_TOKEN_VERSION" envDefault:"1"`
}
