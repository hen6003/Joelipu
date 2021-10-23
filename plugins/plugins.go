package plugins

import (
	"net"
	"net/url"
)

type ServerCfg struct {
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

type GeminiVars struct {
	AbsPath string
	URL *url.URL
	Conn net.Conn
	Cfg ServerCfg
}

type Plugin interface {
	HandleType() string
	HandleGemini(vars GeminiVars) string
}

