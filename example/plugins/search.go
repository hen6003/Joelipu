package main

import (
	"log"
	"os"
	"io/fs"
	"io/ioutil"
	"path/filepath"

	"github.com/sahilm/fuzzy"

	"gmi.hen6003.xyz/joelipu/plugins"
)

type PluginImpl struct{}

func (p PluginImpl) HandleGemini(vars plugins.GeminiVars) string {
	if vars.URL.RawQuery == "" { // If empty query ask for query
		return "10 Search query\r\n"
	}

	// Else create gemini header and title
	data := "20 text/gemini\r\n# Search\n\n## Query: \n>" + string(vars.URL.RawQuery) + "\n\n## Results\n"

	// Read all files in root
	var fileContents []string
	var fileNames []string
	fileSystem := os.DirFS(vars.Cfg.Content.Root)

	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && d.Name()[0] != '.' {
			content, err := ioutil.ReadFile(filepath.Join(vars.Cfg.Content.Root, path))
			if err != nil {
				log.Println(err)
				return nil
			}

			fileContents = append(fileContents, string(content))
			fileNames = append(fileNames, path)
		}

		return nil
	})

	matches := fuzzy.Find(string(vars.URL.RawQuery), fileContents)

	if len(matches) > 0 {
		for _, m := range matches {
			data += "=> " + fileNames[m.Index] + "\n"
		}
	} else {
		data += "No matches found\n"
	}

	data += "\n## Search again \n=> " + vars.URL.Path + " search\n"

	return data
}

func (p PluginImpl) HandleType() string {
	return ".search"
}

var Impl plugins.Plugin = PluginImpl{}

