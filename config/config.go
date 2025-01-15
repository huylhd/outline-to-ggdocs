package config

import "os"

type Config struct {
	OutlineApiKey string
}

func init() {
	AppConfig = &Config{
		OutlineApiKey: os.Getenv("OUTLINE_API_KEY"),
	}
}

var AppConfig *Config
