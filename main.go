package main

import (
	"log"
	"crypto/tls"
	"strconv"
	"plugin"
	"os"
	"net"
	"net/url"

	"gmi.hen6003.xyz/joelipu/plugins"
)

func loadPlugins() Plugin {
	plug, err := plugin.Open("plugin/test.so")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	symPlugin, err := plug.Lookup("Plugin")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	var plugin Plugin
	plugin, ok := symPlugin.(Plugin)
	if !ok {
		log.Println("Unexpected type from module symbol")
		os.Exit(1)
	}

	return plugin
}

func main() {	
	// Load config files
	cfg := loadConfig()
	
	// Load plugins
	plugin := loadPlugins()

	// Load certificate
	cer, err := tls.LoadX509KeyPair(cfg.Certs.Cert, cfg.Certs.Key)
	if err != nil {
		log.Println(err)
		return
	}

	// Setup and start server
	config := &tls.Config{
		Certificates: []tls.Certificate{cer},
		MinVersion:   tls.VersionTLS12,
		ServerName:   cfg.Net.Host,
	}
	ln, err := tls.Listen("tcp", ":"+strconv.Itoa(cfg.Net.Port), config) 
	if err != nil {
		log.Println(err)
		return
	}
	defer ln.Close()

	// Start server loop
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn, cfg, plugin)
	}
}
