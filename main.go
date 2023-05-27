package main

import (
	"os"
	"pioneer-server/io"

	"github.com/phuslu/log"
)

const version = "0.0.1"

func main() {
	os.Setenv("PIONEER_SRV_VERS", version)

	if io.IsCLIMode(os.Args) {
		io.ExecuteCLI(os.Args)
		return
	}

	io.PrintSplashScreen()

	var config = io.GetConfig()
	io.InitLogger(config.Logger)

	log.Info().Msg("Initializing Pioneer Server...")
}
