package main

import (
	"log"
	"os"
	"sync"
)

const Host = "https://www.mixcloud.com"

var wg sync.WaitGroup

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

	for _, link := range config.links {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			mc.Add(link)
		}(link)
	}

	wg.Wait()
}
