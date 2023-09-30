package buildconfig

type Parser interface {
	Parse(path string) (Config, error)
	// CompileConfig templates config file and returns it raw without parsing
	CompileConfig(configPath string) (string, error)
}
