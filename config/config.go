package config

import (
	"fmt"
	"os"
)

type Config struct {
	DbPath string
}

func (c *Config) AutoDiscover() {
	// TODO Retrieve DATA folder depends on tye OS
	home := os.Getenv("HOME")
	if c.DbPath == "" {
		c.DbPath = fmt.Sprintf("%s/.local/share/mytime/mytime.sqlite", home)
	}
}
