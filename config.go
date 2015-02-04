package main

import (
	"bufio"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Config struct {
	Sid           string
	Csrftoken     string
	Playlist_name string
	Links_path    string
	links         []string
}

func (c *Config) Load(path string) {
	confFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	yaml.Unmarshal(confFile, c)
	log.Print(c)

	c.links = readlinks(c.Links_path)
}

func readlinks(path string) []string {
	var links []string
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		links = append(links, strings.TrimSpace(scanner.Text()))
	}
	return links
}
