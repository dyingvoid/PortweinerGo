package config

import (
	"flag"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Env            string         `yaml:"env"`
	Username       string         `yaml:"username"`
	Password       string         `yaml:"password"`
	TerminalConfig TerminalConfig `yaml:"terminal"`
	HttpServer     HTTPServer     `yaml:"http_server"`
}

type TerminalConfig struct {
	UseSudo bool `yaml:"use_sudo"`
}

type HTTPServer struct {
	Port string `yaml:"port"`
}

func MustLoad() *Config {
	configPath := flag.String("config-path", "", "Path to the configuration file.")
	flag.Parse()

	if *configPath == "" {
		*configPath = os.Getenv("CONFIG_PATH")
	}

	if *configPath == "" {
		log.Fatal("Config path is not set. Use -config-path flag or set CONFIG_PATH env variable.")
	}

	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		log.Fatalf("Config file %s does not exist", *configPath)
	}

	data, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("Could not read data from the file: %s", *configPath)
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Could not bind file to config: %v", err)
	}

	return &cfg
}
