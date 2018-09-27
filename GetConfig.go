// Author: Paweł Konopko
// License: MIT

package log

import (
	"io"
	"os"
	"sync"
)

var (
	// mutex will ensure that logs will be printed one by one (not mixed)
	mutex *sync.Mutex = &sync.Mutex{}

	// outputSinks defines where logs should go
	outputSinks []io.Writer = []io.Writer{
		os.Stdout,
	}

	// fieldsOrder defines order of non empty fields in all our logs
	fieldsOrder []string = []string{
		FIELD_TIME,
		FIELD_SESSION,
		FIELD_LEVEL,
		FIELD_MESSAGE,
		CUSTOM_FIELDS,
		FIELD_DATA,
	}

	// logLevel defines minimum log level at witch log messages will be logged
	logLevel LogLevel = LEVEL_INFO

	// customFields defines custom fields that will be added to all logs (just like TIME or LOG_LEVEL)
	customFields Data = Data{}

	// globalSession is parent session for all others sessions
	globalSession Logger = &loggerImpl{
		name: "",
		additionalDataFields: Data{},
	}
)

type Config interface {
	SetOutputSinks (sinks ...io.Writer)
	SetFieldsOrder (keys ...string)
	SetLogLevel (level LogLevel)
	SetCustomFields (fields Data)
}

type Data map[string]interface{}

type LogLevel int
func (level LogLevel) String() string {
	switch level {
	case 0:
		return "fatal"
	case 1:
		return "error"
	case 2:
		return "warning"
	case 3:
		return "info"
	case 4:
		return "debug"
	default:
		return ""
	}
}

func GetConfig() Config {
	return &config{}
}

type config struct {}

func (c *config) SetOutputSinks (sinks ...io.Writer) {
	mutex.Lock()
	defer mutex.Unlock()

	if sinks == nil {
		sinks = []io.Writer{}
	}
	outputSinks = sinks
}

func (c *config) SetFieldsOrder (keys ...string) {
	mutex.Lock()
	defer mutex.Unlock()

	if keys == nil {
		keys = []string{}
	}
	fieldsOrder = keys
}

func (c *config) SetLogLevel (level LogLevel) {
	mutex.Lock()
	defer mutex.Unlock()

	if level <= LEVEL_FATAL {
		logLevel = LEVEL_FATAL
	} else if level >= LEVEL_DEBUG {
		logLevel = LEVEL_DEBUG
	} else {
		logLevel = level
	}
}

func (c *config) SetCustomFields (fields Data) {
	mutex.Lock()
	defer mutex.Unlock()

	if fields == nil {
		fields = Data{}
	}
	customFields = fields
}