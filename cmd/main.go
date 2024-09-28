package main

import (
	"fmt"
	"os"
	"start/internal/config"
	"start/internal/router"

	"github.com/sirupsen/logrus"
	_ "start/docs"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{DisableColors: false})
}

func makeDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
		fmt.Printf("Directory created: %s\n", path)
	} else if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	} else {
		fmt.Printf("Directory already exists: %s\n", path)
	}

	return nil
}

// @title Portweiner
// @version 1.0
// @description API Server for Swarm stack deployment

// @securityDefinitions.basic BasicAuth
// @in header
// @name Authorization
func main() {
	cfg := config.MustLoad()
	logrus.Info("Config loaded.")

	_ = makeDir("./uploads")
	_ = makeDir("./config")
	logrus.Info("Dirs initialized.")

	r := router.SetupRouter(cfg)
	_ = r.Run(":" + cfg.HttpServer.Port)
}
