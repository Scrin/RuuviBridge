package logging

import (
	"os"
	"time"

	"github.com/Scrin/RuuviBridge/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Setup(conf config.Logging) {
	var logger zerolog.Logger

	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano

	switch conf.Level {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	default:
		log.Fatal().Str("configured_level", conf.Level).Msg("Invalid logging level")
	}

	switch conf.Type {
	case "json":
		logger = zerolog.New(os.Stdout)
	case "structured", "simple", "": // simple is not actually supported anymore, but lets keep it for backwards compatibility
		cw := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) { w.Out = os.Stdout })
		logger = zerolog.New(cw)
	default:
		log.Fatal().Str("configured_type", conf.Type).Msg("Invalid logging type")
	}

	ctx := logger.With()
	timestamps := true
	if conf.Timestamps != nil {
		timestamps = *conf.Timestamps
	}
	if timestamps {
		ctx = ctx.Timestamp()
	}
	if conf.WithCaller {
		ctx = ctx.Caller()
	}

	log.Logger = ctx.Logger()
}
