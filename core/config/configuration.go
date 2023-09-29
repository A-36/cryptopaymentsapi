package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Config struct {
	Webserver struct {
		Host    string `json:"host"`
		Debug   bool   `json:"debug"`
		Enabled bool   `json:"enabled"`
	} `json:"webserver"`
}

var (
	Cfg *Config
)

func Load(path string) error {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(f, &Cfg)
	if err != nil {
		return err
	}
	return err
}

func WatchChanges(path string) {
Beginning:
	initialStat, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	for {
		time.Sleep(1 * time.Second)
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			Load(path)
			goto Beginning
		}
	}
}
