package main

import (
	"log"
	"os"
)

const Host = "https://www.mixcloud.com"

func Debug(msg interface{}) {
	if os.Getenv("DEBUG") != "" {
		log.Print(msg)
	}
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config := new(Config)
	config.Load(os.Getenv("HOME") + "/dev/Misc/mixcloud.yml")

	if len(config.links) == 0 {
		log.Print("no links to process")
		return
	}
	Debug(config)

	mc := NewMixcloudPlaylist(config)
	mc.verifyLogin()

	done := make(chan bool)
	for _, link := range config.links {
		go mc.Add(link, done)
	}
	<-done
}
