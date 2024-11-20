package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/1Mochiyuki/gosky/app"
	"github.com/rs/zerolog"
)

var (
	once sync.Once
	log  zerolog.Logger
)

func Get() zerolog.Logger {
	home, err := app.AppHome()
	if err != nil {
		panic(err)
	}
	logFile := fmt.Sprintf("%s/log.log", home)
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if os.IsNotExist(err) {
		panic(err)
	}
	once.Do(func() {
		output := zerolog.ConsoleWriter{
			Out:        file,
			TimeFormat: time.RFC1123,
		}
		output.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}
		output.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("Msg: %s", i)
		}
		output.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("%s: ", i)
		}
		output.FormatFieldValue = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		log = zerolog.New(output).With().Timestamp().Logger()
	})
	return log
}
