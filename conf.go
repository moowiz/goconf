package goconf

import (
	"io"
	"log"
	"strings"
	"time"
)

type config struct {
	modTime time.Time
	info    map[string]interface{}
}

type Config struct {
	name string
}

func (c *Config) Get(key string) interface{} {
	return confMapping[c.name].info[key]
}

var confMapping map[string]Config

func LoadConfig(filename string) (*Config, error) {
	if confMapping == nil {
		confMapping = make(map[string]time.Time)
		confMapping[filename] = time.Now().Add(-100) // Arbitrary time before now to force a reload
		go hotReload()
	}
	return &Config{name:filename}
}

func loadConfig(name string) (map[string]interface{}, error) {
	f, err = os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	split = strings.Split(name, '.')
	format := split[len(split) - 1]
	switch format {
	case "json":
		var rval interface{}
		decoder := json.NewDecoder(f)
		err := decoder.Decode(rval)
		if err != nil {
			return nil, err
		}

		return rval, nil
	default:
		return nil, fmt.Error("Invalid format: %s", format)
	}
}

func hotReload() {
	for {
		for name, conf := range confMapping {
			fi, err := os.Stat(name)
			if err != nil {
				log.Printf("Error reading file %v\n", name)
			}
			if conf.modTime.Before(fi.ModTime()) {
				conf, err := loadConfig(name)
				if err != nil {
					log.Printf("Error reading file %v\n", name)
				}
				confMapping[name] = conf
			}
		}
	}
}