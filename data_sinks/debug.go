package data_sinks

import (
	"encoding/json"
	"fmt"

	"github.com/Scrin/RuuviBridge/parser"
)

func Debug() chan<- parser.Measurement {
	fmt.Println("Starting debug sink")
	measurements := make(chan parser.Measurement)
	go func() {
		for measurement := range measurements {
			data, err := json.Marshal(measurement)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(string(data))
			}
		}
	}()
	return measurements
}
