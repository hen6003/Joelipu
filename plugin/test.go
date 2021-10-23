package main

import (
	"gmi.hen6003.xyz/joelipu/plugins"
)

type PluginImpl struct{}

func (p PluginImpl) HandleGemini(vars plugins.GeminiVars) string {
	return "20 text/plain\r\n" + vars.Cfg.Content.Root
}

func (p PluginImpl) HandleType() string {
	return ".hello"
}

var Impl plugins.Plugin = PluginImpl{}
