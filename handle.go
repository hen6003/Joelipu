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
	"os/exec"
	"mime"

	"gmi.hen6003.xyz/joelipu/plugins"
)

type geminiHeader struct {
	status int
	meta string
}

func handleConnection(conn net.Conn, cfg plugins.ServerCfg, plugin plugins.Plugin) {
	var header geminiHeader
	var content []byte
	var data string 
	var executable bool
	var absPath string

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
	absPath, err = filepath.Abs(filepath.Join(cfg.Content.Root, u.Path))
	if err != nil {
		header.status = 59
		header.meta = "Invalid path"
		goto SEND
	}
	
	// Check for attempting to read files outside of cfg.Root
	if !strings.HasPrefix(absPath, cfg.Content.Root) {
		header.status = 50
		header.meta = "Not allowed"
		goto SEND
	}

	{ // Add index.gmi on if not a folder, in scope due to gotos
		info, err := os.Stat(absPath)
		if err != nil {
			header.status = 51
			header.meta = "File not found"
			goto SEND
		}

		if info.IsDir() {
			// Check for plugins
    	files, err := ioutil.ReadDir(absPath)
    	if err != nil {
				header.status = 51
				header.meta = "File not found"
				goto SEND
    	}

    	for _, f := range files {
				if (f.Name() == plugin.HandleType()) {
					data = plugin.HandleGemini(plugins.GeminiVars{absPath, u, conn, cfg})
					goto WRITE
				}
    	}
 
			absPath = filepath.Join(absPath, cfg.Content.Index)
		}
		
		// Check again that index.gmi exists, and get the file mode
		info, err = os.Stat(absPath)
		if err != nil {
			header.status = 51
			header.meta = "File not found"
			goto SEND
		}

		executable = info.Mode().Perm() & 0111 != 0
	}

	if !executable {
		// Finally read file
		content, err = ioutil.ReadFile(absPath)
		if err != nil {
			header.status = 51
			header.meta = "File not found"
			goto SEND
		}

		// Get mime type from file
		header.meta = mime.TypeByExtension(filepath.Ext(absPath))
		header.status = 20	
	} else {
		cmd := exec.Command(absPath)

		// Create envvars
		cmd.Env = append(os.Environ(),
			// Request info
			"SCRIPT_FILENAME="+absPath,
			"SCRIPT_NAME="+u.Path,
			"REQUEST_URI="+msg,
			"QUERY_STRING="+u.RawQuery,
			// Remote info
			"REMOTE_ADDR="+conn.RemoteAddr().String(),
			"REMOTE_HOST="+conn.RemoteAddr().String(),
			// Server info
			"GATEWAY_INTERFACE=CGI/1.1",
			"DOCUMENT_ROOT="+cfg.Content.Root,
			"SERVER_NAME="+cfg.Net.Host,
			"SERVER_PORT="+strconv.Itoa(cfg.Net.Port),
			"SERVER_SOFTWARE=Joelipu",
		)

		// Run command and get output
		out, err := cmd.Output()
		if err != nil {
			header.status = 42
			header.meta = "CGI Failed"
			goto SEND
		}

		// Finally send data, command output uses its own
		n, err := conn.Write(out)
		if err != nil {
			log.Println(n, err)
		}
		return
	}

SEND:
	if (header.status == 20) { // If header.status 20 add response
		data = strconv.Itoa(header.status) + " " + header.meta + "\r\n" + string(content)
	} else { // else omit it
		data = strconv.Itoa(header.status) + " " + header.meta + "\r\n"
	}

WRITE:
	// Finally send data
	n, err := conn.Write([]byte(data))
	if err != nil {
		log.Println(n, err)
		return
	}
}
