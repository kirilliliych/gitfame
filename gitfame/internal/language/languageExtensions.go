package language

import (
	"encoding/json"
	"os"
	"strings"
)

type Language struct {
	Name       string
	Extensions []string
}

var langMap = make(map[string][]string)
var configRead bool = false

func readConfig() {
	configRead = true
	languageExtensions, _ := os.ReadFile("./../../configs/language_extensions.json")
	var languages []Language
	_ = json.Unmarshal(languageExtensions, &languages)
	for _, language := range languages {
		langMap[strings.ToLower(language.Name)] = language.Extensions
	}
}

func GetLanguageExtensions(language string) []string {
	if !configRead {
		readConfig()
	}
	return langMap[strings.ToLower(language)]
}
