package main

import (
	"log"
)

const Host = "http://www.mixcloud.com"

func main() {

	config := new(Config)
	config.Load("/home/nenadpet/dev/Misc/mixcloud.yml")

	if len(config.links) == 0 {
		log.Print("no links to process")
		return
	}
	log.Print(config)

	mc := NewMixcloudPlaylist(config)
	mc.verifyLogin()

	for _, link := range config.links {
		mc.Add(link)
	}
}
