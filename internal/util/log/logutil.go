// The logutil provides very simple structured logging with text color that matches the log level
package logutil

import (
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	*log.Logger
}

type LogLevelId int

type LogLevel struct {
	name  string
	level LogLevelId
	color string
}

const (
	INFO_LEVEL  LogLevelId = 1
	ERROR_LEVEL LogLevelId = 2
	WARN_LEVEL  LogLevelId = 3

	colorRed    string = "\033[31m"
	colorGreen  string = "\033[32m"
	colorYellow string = "\033[33m"
	colorReset  string = "\033[0m"

	// These constants are used to grab the file name from the stack when making a logger call.
	// If the code to create or get a logger is changed, these should be examined to ensure they are still the correct values.

	loggerCreateSkipLevel int = 7
	loggerGetSkipLevel    int = 3
)

var logLevels = make(map[LogLevelId]LogLevel)

var cmlogger *Logger
var once sync.Once

func Log(levelId LogLevelId, message string) {
	logger := Get()
	logLevel := logLevels[levelId]
	logger.Println(string(logLevel.color), logLevel.name, string(colorReset), message)
}

func LogError(err error) {
	Log(ERROR_LEVEL, err.Error())
}

func Get() *Logger {
	once.Do(func() {
		cmlogger = createLogger()
		initLevels()
	})

	cmlogger.SetPrefix(loggerPrefix(loggerGetSkipLevel))

	return cmlogger
}

func createLogger() *Logger {
	iow := io.Writer(os.Stdout)

	return &Logger{
		Logger: log.New(iow, loggerPrefix(loggerCreateSkipLevel), log.Lmsgprefix),
	}
}

func initLevels() {
	logLevels[INFO_LEVEL] = LogLevel{
		name:  "INFO",
		level: INFO_LEVEL,
		color: colorGreen,
	}

	logLevels[ERROR_LEVEL] = LogLevel{
		name:  "ERROR",
		level: ERROR_LEVEL,
		color: colorRed,
	}

	logLevels[WARN_LEVEL] = LogLevel{
		name:  "WARN",
		level: WARN_LEVEL,
		color: colorYellow,
	}
}

func loggerPrefix(skip int) string {
	dateTime := time.Now().Local().UTC().Format("2006-01-02 15:04:05")

	_, file, line, _ := runtime.Caller(skip)

	directories := strings.Split(file, "/")

	calledFrom := directories[len(directories)-1]

	return dateTime + " " + calledFrom + ":" + strconv.Itoa(line)
}
