package config

type Config struct {
	Secrets []Secret
}

type Secret struct {
	ID   string
	Path string
}
