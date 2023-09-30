package config

type Config struct {
	Secrets []Secret `json:"secrets"`
}

type Secret struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}
