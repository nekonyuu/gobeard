package util

import (
	"io/ioutil"
	"time"

	"github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	CheckInterval   time.Duration  `yaml:"check_interval"`
	EventsQueueSize int            `yaml:"events_queue_size"`
	Api             ApiConfig      `yaml:"api"`
	MongoDb         MongoDbConfig  `yaml:"mongodb"`
	Torrents        TorrentsConfig `yaml:"torrents"`
	Slack           SlackConfig    `yaml:"slack"`
}

type MongoDbConfig struct {
	Host string `yaml:"host"`
}

type ApiConfig struct {
	Addr string `yaml:"addr"`
}

type TorrentsConfig struct {
	Transmission TransmissionConfig `yaml:"transmission"`
	Quality      []string           `yaml:"quality"`
	WatchDir     string             `yaml:"watch_dir"`
}

type TransmissionConfig struct {
	Endpoint    string `yaml:"endpoint"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	DownloadDir string `yaml:"download_dir"`
}

type SlackConfig struct {
	WebhookUrl string `yaml:"webhook_url"`
	Channel    string `yaml:"channel"`
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
