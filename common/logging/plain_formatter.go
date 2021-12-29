package logging

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type PlainFormatterWithTsWithCaller struct {
}
type PlainFormatterWithoutTsWithCaller struct {
}
type PlainFormatterWithTsWithoutCaller struct {
}
type PlainFormatterWithoutTsWithoutCaller struct {
}

func (f *PlainFormatterWithTsWithCaller) Format(entry *log.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("[%s] [%s] [%s] %s\n", entry.Time.Format("2006-1-2 15:04:05"), entry.Level, entry.Caller.File, entry.Message)), nil
}
func (f *PlainFormatterWithoutTsWithCaller) Format(entry *log.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("[%s] [%s] %s\n", entry.Level, entry.Caller.File, entry.Message)), nil
}

func (f *PlainFormatterWithTsWithoutCaller) Format(entry *log.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("[%s] [%s] %s\n", entry.Time.Format("2006-1-2 15:04:05"), entry.Level, entry.Message)), nil
}
func (f *PlainFormatterWithoutTsWithoutCaller) Format(entry *log.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("[%s] %s\n", entry.Level, entry.Message)), nil
}
