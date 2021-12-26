package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Scrin/RuuviBridge/common/version"
	"github.com/Scrin/RuuviBridge/config"
	"github.com/Scrin/RuuviBridge/processor"
)

func main() {
	configPath := flag.String("config", "./config.yml", "The path to the configuration")
	versionFlag := flag.Bool("version", false, "Prints the version and exits")
	flag.Parse()

	if *versionFlag {
		version.Print()
		return
	}
	conf, err := config.ReadConfig(*configPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	processor.Run(conf)
}
