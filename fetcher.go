package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type MixcloudPlaylist struct {
	config *Config
}

var httpClient = &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
	// same path, 301 redirect to https
	if len(via) == 1 && req.URL.Path == via[0].URL.Path {
		Debug("allow redirect")
		return nil
	} else {
		Debug("skip redirect")
		return errors.New("skip redirect")
	}
},
}

func NewMixcloudPlaylist(c *Config) *MixcloudPlaylist {
	return &MixcloudPlaylist{config: c}
}

func (m *MixcloudPlaylist) verifyLogin() {
	req, err := http.NewRequest("GET", Host, nil)
	req.Header.Add("Cookie", m.cookies())
	Debug(fmt.Sprintf("req:%+v", req))
	if err != nil {
		log.Fatal(err)
	}
	resp, err := httpClient.Do(req)
	//if err != nil {
	//	log.Fatal(err)
	//}
	if resp.StatusCode != 200 {
		log.Fatal(resp)
	}
	Debug(fmt.Sprintf("resp:%+v", resp))
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	if !strings.Contains(string(body), `_loggedIn": true`) {
		log.Fatal("Invalid session id!")
	}
}

func (m *MixcloudPlaylist) Add(link string) {
	path := m.parsePath(link)
	m.post(path)
}

func (m *MixcloudPlaylist) cookies() string {
	return fmt.Sprintf("s=%s; csrftoken=%s", m.config.Sid, m.config.Csrftoken)
}

func (m *MixcloudPlaylist) parsePath(link string) string {
	var path string

	Debug(fmt.Sprintf("link:%+v", link))
	req, err := http.NewRequest("GET", link, nil)
	Debug(fmt.Sprintf("req:%+v", req))
	if err != nil {
		log.Println(link)
		log.Fatal(err)
	}

	resp, err := httpClient.Do(req)
	defer resp.Body.Close()
	Debug(fmt.Sprintf("resp:%+v", resp))

	switch resp.StatusCode {
	case 200:
		url, err := url.Parse(link)
		if err != nil {
			log.Println(link)
			log.Fatal(err)
		}
		path = url.Path
	case 302:
		url, err := resp.Location()
		if err != nil {
			log.Println(link)
			log.Fatal(err)
		}
		path = url.Path
	default:
		log.Println(link)
		log.Fatal(resp)
	}
	Debug(fmt.Sprintf("path:%+v", path))

	return path
}

func (m *MixcloudPlaylist) post(path string) {
	fullPath := fmt.Sprintf("/playlists%sadd-to-collection/", path)

	body := []byte(fmt.Sprintf("action=add&playlist_slug=%s", m.config.Playlist_name))
	req, err := http.NewRequest("POST", Host+fullPath, bytes.NewBuffer(body))
	if err != nil {
		log.Println(path)
		log.Fatal(err)
	}
	Debug(fmt.Sprintf("req:%+v", req))

	req.Header.Add("Cookie", m.cookies())
	req.Header.Add("Origin", Host)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Referer", Host+path)
	req.Header.Add("X-CSRFToken", m.config.Csrftoken)
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	Debug(fmt.Sprintf("req:%+v", req))

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println(req)
		log.Println(path)
		log.Fatal(err)
	}
	defer resp.Body.Close()
	Debug(fmt.Sprintf("resp:%+v", resp))

	if resp.StatusCode != 200 {
		log.Println(req)
		log.Println(path)
		log.Fatal(resp)
	}
}
