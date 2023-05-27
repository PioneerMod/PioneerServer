package io

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

func IsCLIMode(args []string) bool {
	return len(args) > 1 && args[1] == "cli"
}

func ExecuteCLI(args []string) {
	if len(args) < 3 {
		printHelp()
		return
	}

	switch args[2] {
	case "--help", "-h":
		printHelp()
	case "--version", "-v":
		printVersion()
	case "--generate-config", "--gen-config", "-gc":
		generateConfig()
	default:
		printHelp()
	}
}

func printHelp() {
	fmt.Println("Usage: pioneer-server cli --[command]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  --help, -h")
	fmt.Println("  --version, -v")
	fmt.Println("  --generate-config, --gen-config, -gc")
}

func printVersion() {
	fmt.Println("Pioneer Server v" + os.Getenv("PIONEER_SRV_VERS") + " by DasDarki & Pythagorion")
}

func generateConfig() {
	fmt.Println("Should the server run only on the local machine? (y/n)")
	isLocal := readBoolInput()
	fmt.Println("")

	fmt.Println("Which port should the server listen on?")
	var port int
	fmt.Scanln(&port)
	fmt.Println("")

	fmt.Println("Should the server be password protected? (y/n)")
	isPasswordProtected := readBoolInput()
	fmt.Println("")

	var password *string
	if isPasswordProtected {
		fmt.Println("Enter the password:")
		var passwordString string
		fmt.Scanln(&passwordString)
		password = &passwordString
	}

	fmt.Println("")
	fmt.Println("Generating config...")

	var config = Config{
		Server: ServerConfig{
			IsLocal:  isLocal,
			Port:     port,
			Password: password,
		},
		Logger: LoggerConfig{
			Level:          "info",
			EnableColorize: true,
			TimeFormat:     "15:04:05.000",
		},
	}

	f, err := os.Create("config.toml")
	if err != nil {
		panic(err)
	}

	encoder := toml.NewEncoder(f)
	err = encoder.Encode(config)
	if err != nil {
		panic(err)
	}

	fmt.Println("Config generated!")
}

func readBoolInput() bool {
	var input string
	fmt.Scanln(&input)
	return input == "y"
}
