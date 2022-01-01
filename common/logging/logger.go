package logging

import (
	"sort"

	"github.com/Scrin/RuuviBridge/config"
	log "github.com/sirupsen/logrus"
)

func Setup(conf config.Logging) {
	timestamps := true
	if conf.Timestamps != nil {
		timestamps = *conf.Timestamps
	}

	log.SetReportCaller(conf.WithCaller)

	if conf.WithCaller {
		if timestamps {
			log.SetFormatter(new(PlainFormatterWithTsWithCaller))
		} else {
			log.SetFormatter(new(PlainFormatterWithoutTsWithCaller))
		}
	} else {
		if timestamps {
			log.SetFormatter(new(PlainFormatterWithTsWithoutCaller))
		} else {
			log.SetFormatter(new(PlainFormatterWithoutTsWithoutCaller))
		}
	}

	switch conf.Type {
	case "structured":
		log.SetFormatter(&log.TextFormatter{
			DisableTimestamp: !timestamps,
			SortingFunc:      sortFN,
		})
	case "json":
		log.SetFormatter(&log.JSONFormatter{
			DisableTimestamp: !timestamps,
		})
	case "simple":
	case "":
	default:
		log.Fatal("Invalid logging type: ", conf.Type)
	}

	switch conf.Level {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	case "":
		log.SetLevel(log.InfoLevel)
	default:
		log.Fatal("Invalid logging level: ", conf.Level)
	}
}

func sortFN(keys []string) {
	sort.Slice(keys, func(i, j int) bool {
		switch keys[i] {
		case "time":
			return true
		case "level":
			return keys[j] != "time"
		case "msg":
			return keys[j] != "time" && keys[j] != "level"
		case "error":
			return keys[j] == "file" || keys[j] == "func"
		case "func":
			return keys[j] == "file"
		case "file":
			return false
		}
		switch keys[j] {
		case "time":
			return false
		case "level":
			return keys[j] == "time"
		case "msg":
			return keys[i] == "level"
		case "error":
			return keys[i] != "file" && keys[i] != "func"
		case "func":
			return keys[i] != "file"
		case "file":
			return true
		}
		return keys[i] < keys[j]
	})
}
