package utils

import (
	"flag"
	"fmt"
	"github.com/siddontang/go-log/log"
	"gopkg.in/yaml.v3"
	config "iot-x/conf"
	"os"
)

var release string

func init() {
	flag.StringVar(&release, "release", "local", "release model, optional local/dev/prod")
}

func loadConfig(filePath string) (*config.Config, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var conf config.Config
	err = yaml.Unmarshal(content, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func InitConf() *config.Config {
	flag.Parse()

	s := fmt.Sprintf("etc/application-%s.yaml", release)

	conf, err := loadConfig(s)

	if err != nil {
		log.Fatal(err)
	}

	conf.Server.IP = getLocalIP()

	return conf
}
