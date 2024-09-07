package main

import (
	"github.com/passionde/user-segmentation-service/internal/app"
)

const configPath = "config/config.yaml"

func main() {
	app.Run(configPath)
}
