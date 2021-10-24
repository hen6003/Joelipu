package main

import (
	"github.com/hen6003/joelipu/plugins"
)

type PluginImpl struct{}

func (p PluginImpl) HandleGemini(vars plugins.GeminiVars) string {
	return "10 Hello World\r\n"
}

func (p PluginImpl) HandlePath() string {
	return "hello"
}

var Impl plugins.Plugin = PluginImpl{}
