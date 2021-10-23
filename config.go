package main

import (
	"log"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"gmi.hen6003.xyz/joelipu/plugins"
)

func loadConfig() plugins.ServerCfg {
	var cfg plugins.ServerCfg

	// Defaults
	cfg.Net.Host = "localhost"
	cfg.Net.Port = 1965
	cfg.Content.Root = "root"
	cfg.Content.Index = "index.gmi"
	cfg.Content.Plugins = "plugins"

	// Read config file
	if _, err := toml.DecodeFile("config.toml", &cfg); err != nil {
		log.Println(err)
		return cfg
	}

	var err error
	cfg.Content.Root, err = filepath.Abs(cfg.Content.Root)
	if err != nil {
		log.Println(err)
		return cfg
	}

	return cfg
}
