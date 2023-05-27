package io

import (
	"fmt"
	"os"

	"github.com/phuslu/log"
)

func InitLogger(config LoggerConfig) {
	log.DefaultLogger.Writer = &log.ConsoleWriter{
		ColorOutput:    true,
		QuoteString:    true,
		EndWithMessage: true,
	}

	log.DefaultLogger.Level = log.ParseLevel(config.Level)
	log.DefaultLogger.TimeFormat = config.TimeFormat
}

func PrintSplashScreen() {
	fmt.Println("")
	fmt.Println(`    ____  _                            __  ___          __
   / __ \(_)___  ____  ___  ___  _____/  |/  /___  ____/ /
  / /_/ / / __ \/ __ \/ _ \/ _ \/ ___/ /|_/ / __ \/ __  / 
 / ____/ / /_/ / / / /  __/  __/ /  / /  / / /_/ / /_/ /  
/_/   /_/\____/_/ /_/\___/\___/_/  /_/  /_/\____/\__,_/   
                                                          `)
	fmt.Println("")
	fmt.Println("Pioneer Server v" + os.Getenv("PIONEER_SRV_VERS") + " by DasDarki & Pythagorion")
	fmt.Println("")
	fmt.Println("")
}
