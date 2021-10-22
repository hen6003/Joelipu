# Joelipu
Gemini server written for my capsule (gemini://gmi.hen6003.xyz/)

## Features
* Static file hosting
* CGI file hosting

## Config
The server is configured in TOML from a config.toml file in the directory the program is ran.

### Content
* root: The directory (relative) holding the sites files (defualts to root)
* index: What file to use if a directory is requested (defualts to index.gmi)

### Net
* port: Port to host on (defualts to 1965)
* host: Server name (defaults to localhost)

### Certs
* cert: Path to certificate
* key: Path to private key

### Example
```Example config
[Content]
root = "root"
index = "index.gmi"

[Net]
port = 1965
host = "server.name"

[Certs]
cert = "certs/server.crt"
key  = "certs/server.key"
```

## CGI
If a file is executable the server will execute the file and serve the output. The following enviroment variables will be set:

### Request info
* SCRIPT_FILENAME: The absolute path to the script
* SCRIPT_NAME: The relative path to the script from DOCUMENT_ROOT
* REQUEST_URI: The complete URI that was requested
* QUERY_STRING: The query string in the URI

### Remote info
* REMOTE_ADDR: The IP address of the visitor
* REMOTE_HOST: Their IP address again (For compatability with the spec)

### Server info
* GATEWAY_INTERFACE: "CGI/1.1" (For compatability with the spec)
* DOCUMENT_ROOT: The root of the hosted content
* SERVER_NAME: The servers hostname (e.g. gmi.hen6003.xyz)
* SERVER_PORT: The port the server is hosting on
* SERVER_SOFTWARE: "Joelipu" (For compatability with the spec)

## Setup certificates
Make a folder called 'certs', and use openssl to create the certificates
```Example command
$ mkdir certs
$ openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:4096 -keyout certs/server.key -out certs/server.crt
```

## Name
"jo e lipu" is Toki Pona for "Having paper". "pana" (giving) is probally a better word than "jo", however "panaelipu" looks worse ¯\\\_(ツ)\_/¯.
