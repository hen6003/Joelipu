package main

import (
	"gmi.hen6003.xyz/joelipu/plugins"
)

type plugin bool

func (p plugin) HandleGemini(vars ServerCfg) string {
	return "20 text/plain\r\n" + vars.Content.Root
}

func (p plugin) HandleType() string {
	return ".hello"
}

var Plugin plugin
