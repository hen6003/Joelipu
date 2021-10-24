// Regex search plugin

package main

import (
	"log"
	"os"
	"io/fs"
	"io/ioutil"
	"path/filepath"

	"regexp"

	"gmi.hen6003.xyz/joelipu/plugins"
)

type PluginImpl struct{}

func (p PluginImpl) HandleGemini(vars plugins.GeminiVars) string {
	var fileSystem fs.FS
	var fileContents []string
	var fileNames []string
	matches := false

	if vars.URL.RawQuery == "" { // If empty query ask for query
		return "10 Search query\r\n"
	}

	// Else create gemini header and title
	data := "20 text/gemini\r\n# Regex search\n\n## Query: \n>" + string(vars.URL.RawQuery) + "\n\n"

	// Compile regex
	regex, err := regexp.Compile(string(vars.URL.RawQuery))
	if err != nil { // If failed return so
		data += "## Failed\nRegex failed to compile"
		goto END
	}

	// Else show results
	data += "## Results\n"

	// Read all files in root
	fileSystem = os.DirFS(vars.Cfg.Content.Root)

	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && d.Name()[0] != '.' {
			content, err := ioutil.ReadFile(filepath.Join(vars.Cfg.Content.Root, path))
			if err != nil {
				log.Println(err)
				return nil
			}

			fileContents = append(fileContents, string(content))
			fileNames = append(fileNames, path)
		} else if d.Name()[0] == '.' && len(d.Name()) > 1 {
			return fs.SkipDir // Ignore hidden directorys
		}

		return nil
	})

	for i, c := range fileContents {
		if regex.MatchString(c) {
			data += "=> " + fileNames[i] + "\n"
			matches = true
		}
	}

	if !matches {
		data += "No matches found\n"
	}

END:
	data += "\n## Search again \n=> " + vars.URL.Path + " search\n"

	return data
}

func (p PluginImpl) HandlePath() string {
	return "search"
}

var Impl plugins.Plugin = PluginImpl{}

