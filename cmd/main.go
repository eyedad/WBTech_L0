package main

import (
	"example.com/m/v2/config"
	"example.com/m/v2/internal/app"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	app.Run(cfg)
}
