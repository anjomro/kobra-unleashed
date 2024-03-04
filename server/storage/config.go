package storage

type Config struct {
	Folder string
}

var ConfigDefault = Config{
	Folder: "./.sessions",
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]

	if cfg.Folder == "" {
		cfg.Folder = ConfigDefault.Folder
	}

	return cfg
}
