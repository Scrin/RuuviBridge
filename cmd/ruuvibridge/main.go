package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"github.com/Scrin/RuuviBridge/common/logging"
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
	logging.Setup(conf.Logging) // logging should be set up with logging config before logging a possible error in the config, weird, I know
	if err != nil {
		log.Panic(err)
	}
	processor.Run(conf)
}
