package main

import (
	"log"
	"crypto/tls"
	"strconv"
	"plugin"
	"os"
	"io/ioutil"
	"path/filepath"

	"github.com/hen6003/joelipu/plugins"
)

func loadPlugins(path string) []plugins.Plugin {
	var pluginlist []plugins.Plugin

  files, err := ioutil.ReadDir(path)
  if err != nil {
		log.Println(err)
  }

  for _, f := range files {
		plug, err := plugin.Open(filepath.Join(path,f.Name()))
		if err != nil {
			log.Println(err, "... skipping")
			continue
		}

		symPlugin, err := plug.Lookup("Impl")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		newplugin, ok := symPlugin.(*plugins.Plugin)
		if !ok {
			log.Println("Unexpected type")
			os.Exit(1)
		}

		pluginlist = append(pluginlist, *newplugin)
	}

	return pluginlist
}

func main() {	
	// Load config files
	cfg := loadConfig()
	
	// Load plugins
	pluginlist := loadPlugins(cfg.Content.Plugins)

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
		go handleConnection(conn, cfg, pluginlist)
	}
}
