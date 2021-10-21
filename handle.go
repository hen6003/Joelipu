package main

import (
	"log"
	"net"
	"net/url"
	"bufio"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"strconv"
	"os"
	"mime"
)

type geminiHeader struct {
	status int
	meta string
}

func handleConnection(conn net.Conn, cfg serverCfg) {
	var header geminiHeader
	var content []byte

	defer conn.Close()
	r := bufio.NewReader(conn)

	msg, err := r.ReadString('\r')
	if err != nil {
		if err != io.EOF { // EOF doesn't need to be logged
			log.Println(err)
		}
		return
	}

	msg = msg[:len(msg)-1] // Remove '\r'

	// Parse url
	u, err := url.Parse(msg)
	if err != nil {
		log.Println(err)
		header.status = 59
		header.meta = "Unable to read url"
		goto SEND
	}

	// Check for proxying
	if u.Scheme != "gemini" ||
     u.Host != cfg.Net.Host {
		header.status = 53
		header.meta = "Proxying is not supported"
		goto SEND
	}

	// Create absolute path
	u.Path, err = filepath.Abs(filepath.Join(cfg.Content.Root, u.Path))
	if err != nil {
		header.status = 59
		header.meta = "Invalid path"
		goto SEND
	}
	
	// Check for attempting to read files outside of cfg.Root
	if !strings.HasPrefix(u.Path, cfg.Content.Root) {
		header.status = 50
		header.meta = "Not allowed"
		goto SEND
	}

	{ // Add index.gmi on if not a folder, in scope due to gotos
		info, err := os.Stat(u.Path)
		if err != nil {
			header.status = 51
			header.meta = "File not found"
			goto SEND
		}

		if info.IsDir() {
			u.Path = filepath.Join(u.Path, cfg.Content.Index)
		}
	}

	// Finally read file
	content, err = ioutil.ReadFile(u.Path)
	if err != nil {
		header.status = 51
		header.meta = "File not found"
		goto SEND
	}

	// Get mime type from file
	header.meta = mime.TypeByExtension(filepath.Ext(u.Path))
	header.status = 20

SEND:
	var data string 

	if (header.status == 20) { // If header.status 20 add response
		data = strconv.Itoa(header.status) + " " + header.meta + "\r\n" + string(content)
	} else { // else omit it
		data = strconv.Itoa(header.status) + " " + header.meta + "\r\n"
	}

	// Finally send data
	n, err := conn.Write([]byte(data))
	if err != nil {
		log.Println(n, err)
		return
	}	
}
