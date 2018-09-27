package log

func Session(name string) Logger {
	return globalSession.Session(name)
}

func Fatal(message string, data ...Data) {
	globalSession.Fatal(message, data...)
}

func Error(message string, data ...Data) {
	globalSession.Error(message, data...)
}

func Warning(message string, data ...Data) {
	globalSession.Warning(message, data...)
}

func Info(message string, data ...Data) {
	globalSession.Info(message, data...)
}

func Debug(message string, data ...Data) {
	globalSession.Debug(message, data...)
}