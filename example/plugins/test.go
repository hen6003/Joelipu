package main

import (
	"gmi.hen6003.xyz/joelipu/plugins"
)

type PluginImpl struct{}

func (p PluginImpl) HandleGemini(vars plugins.GeminiVars) string {
	return "10 Hello World\r\n"
}

func (p PluginImpl) HandlePath() string {
	return "hello"
}

var Impl plugins.Plugin = PluginImpl{}
