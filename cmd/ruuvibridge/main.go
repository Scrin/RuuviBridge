package main

import (
	"flag"

	"github.com/rs/zerolog/log"

	"github.com/Scrin/RuuviBridge/common/logging"
	"github.com/Scrin/RuuviBridge/common/version"
	"github.com/Scrin/RuuviBridge/config"
	"github.com/Scrin/RuuviBridge/processor"
)

func main() {
	configPath := flag.String("config", "./config.yml", "The path to the configuration")
	strictConfig := flag.Bool("strict-config", false, "Use strict parsing for the config file; will throw errors for invalid fields")
	versionFlag := flag.Bool("version", false, "Prints the version and exits")
	flag.Parse()

	if *versionFlag {
		version.Print()
		return
	}

	conf, err := config.ReadConfig(*configPath, *strictConfig)
	logging.Setup(conf.Logging) // logging should be set up with logging config before logging a possible error in the config, weird, I know
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}
	log.Debug().Str("configfile", *configPath).Msg("Config loaded")
	processor.Run(conf)
}
