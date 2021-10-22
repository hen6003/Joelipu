package main

import (
	"log"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type serverCfg struct {
	Net netInfo
	Certs certsInfo
	Content contentInfo
}

type certsInfo struct {
	Path string
	Cert string
	Key string
}

type netInfo struct {
	Host string
	Port int
}

type contentInfo struct {
	Root string
	Index string
}

func loadConfig() serverCfg {
	var cfg serverCfg

	// Defaults
	cfg.Net.Host = "localhost"
	cfg.Net.Port = 1965
	cfg.Content.Root = "root"
	cfg.Content.Index = "index.gmi"

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
