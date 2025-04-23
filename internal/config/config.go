package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBUrl string
}

func Load() Config {
	home := os.Getenv("HOME")
	return Config{
		DBUrl: fmt.Sprintf("file://%s/.local/share/mytime/mytime.sqlite", home),
	}
}
