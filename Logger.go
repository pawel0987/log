// Author: Paweł Konopko
// License: MIT

package log

import "time"

const (
	LEVEL_FATAL   = LogLevel(0)
	LEVEL_ERROR   = LogLevel(1)
	LEVEL_WARNING = LogLevel(2)
	LEVEL_INFO    = LogLevel(3)
	LEVEL_DEBUG   = LogLevel(4)

	FIELD_TIME    = "time"
	FIELD_SESSION = "session"
	FIELD_LEVEL   = "level"
	FIELD_MESSAGE = "message"
	FIELD_DATA    = "data"
	CUSTOM_FIELDS = "customFields"
)

type Logger interface {
	// Fatal will log message on fatal level.
	Fatal (message string, d ...Data)

	// Error will log message on error level.
	Error (message string, d ...Data)

	// Warning will log message on warning level.
	Warning (message string, d ...Data)

	// Info will log message on info level.
	Info (message string, d ...Data)

	// Debug will log message on debug level.
	Debug (message string, d ...Data)

	// AddDataField will add specified field to logger session.
	// All future logging will include this fields.
	AddDataField (key string, field interface{})

	// Session will create sub session from current logger session.
	// Session name will be <CURRENT_SESSION_NAME>.<SUB_SESSION_NAME>
	// or just <SUB_SESSION_NAME> when current session name was empty.
	Session (subname string) Logger
}

type loggerImpl struct {
	name                 string
	additionalDataFields Data
}

func (lg *loggerImpl) AddDataField(key string, field interface{}) {
	mutex.Lock()
	defer mutex.Unlock()

	lg.additionalDataFields[key] = field
}

func (lg *loggerImpl) Session(subname string) Logger {
	previousDataFields := Data{}
	for key, value := range lg.additionalDataFields {
		previousDataFields[key] = value
	}

	return &loggerImpl{
		name: func() string {
			if len(lg.name) != 0 {
				return lg.name + "." + subname
			} else {
				return subname
			}
		}(),
		additionalDataFields: previousDataFields,
	}
}

func (lg *loggerImpl) Fatal(msg string, d ...Data) {
	lg.sendLog(LEVEL_FATAL, msg, d)
}

func (lg *loggerImpl) Error(msg string, d ...Data) {
	if logLevel >= LEVEL_ERROR {
		lg.sendLog(LEVEL_ERROR, msg, d)
	}
}

func (lg *loggerImpl) Warning(msg string, d ...Data) {
	if logLevel >= LEVEL_WARNING {
		lg.sendLog(LEVEL_WARNING, msg, d)
	}
}

func (lg *loggerImpl) Info(msg string, d ...Data) {
	if logLevel >= LEVEL_INFO {
		lg.sendLog(LEVEL_INFO, msg, d)
	}
}

func (lg *loggerImpl) Debug(msg string, d ...Data) {
	if logLevel >= LEVEL_DEBUG {
		lg.sendLog(LEVEL_DEBUG, msg, d)
	}
}

func (lg *loggerImpl) sendLog (level LogLevel, message string, data []Data) {
	keysOrder := []string{}
	fields := map[string]interface{}{
		FIELD_TIME: time.Now().String(),
		FIELD_LEVEL: level.String(),
		FIELD_MESSAGE: message,
		FIELD_SESSION: lg.name,
	}

	for _, key := range fieldsOrder {
		if key == CUSTOM_FIELDS {
			for customKey, customValue := range customFields {
				keysOrder = append(keysOrder, customKey)
				fields[customKey] = customValue
			}
		} else if key == FIELD_DATA {
			if len(data) != 0  || len(lg.additionalDataFields) != 0 {
				keysOrder = append(keysOrder, FIELD_DATA)
				fields[FIELD_DATA] = joinData(append(data, lg.additionalDataFields))
			}
		} else {
			keysOrder = append(keysOrder, key)
		}
	}

	if len(message) == 0 {
		delete(fields, FIELD_MESSAGE)
		for id, v := range keysOrder {
			if v == FIELD_MESSAGE {
				keysOrder = append(keysOrder[0:id], keysOrder[id+1:]...)
				break
			}
		}
	}

	log := formatJsonMessage(keysOrder, fields)
	mutex.Lock()
	defer mutex.Unlock()
	for _, sink := range outputSinks {
		sink.Write([]byte(log))
	}
}
