package main

import (
	"log"
	"crypto/tls"
	"strconv"
)

func main() {
	// Load config files
	cfg := loadConfig()

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
		go handleConnection(conn, cfg)
	}
}
