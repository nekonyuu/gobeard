package util

import (
	"io/ioutil"
	"time"

	"github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	CheckInterval   time.Duration `yaml:"check_interval"`
	EventsQueueSize int           `yaml:"events_queue_size"`
	Api             struct {
		Addr string `yaml:"addr"`
	} `yaml:"api"`
	MongoDb struct {
		Host string `yaml:"host"`
	} `yaml:"mongodb"`
	Torrents struct {
		Transmission struct {
			Endpoint    string `yaml:"endpoint"`
			Username    string `yaml:"username"`
			Password    string `yaml:"password"`
			DownloadDir string `yaml:"download_dir"`
		} `yaml:"transmission"`
		Quality  []string `yaml:"quality"`
		WatchDir string   `yaml:"watch_dir"`
	} `yaml:"torrents"`
	Slack struct {
		WebhookUrl string `yaml:"webhook_url"`
		Channel    string `yaml:"channel"`
	} `yaml:"slack"`
	Actions     []string `yaml:"actions"`
	Downloaders []string `yaml:"downloaders"`
}

var config Config

func GetConfig() Config {
	return config
}

func SetConfig(configPath string) {
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		logrus.Fatalf("error parsing configuration file: %s", err)
	}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		logrus.Fatalf("error parsing configuration file: %s", err)
	}
}
