package logger

type ConsoleLogger struct {
	Level LogLv
}

func NewConsoleLogger(level string) *ConsoleLogger {
	l := &ConsoleLogger{
		Level: str2Lv(level),
	}

	return l
}

func (l *ConsoleLogger) enable(logLv LogLv) bool {
	return l.Level <= logLv
}

func (l *ConsoleLogger) Debug(str string, a ...any) {
	lv := DEBUG
	if !l.enable(lv) {
		return
	}
	logPrint(lv, str, a...)
}

func (l *ConsoleLogger) Trace(str string, a ...any) {
	lv := TRACE
	if !l.enable(lv) {
		return
	}
	logPrint(lv, str, a...)
}

func (l *ConsoleLogger) Info(str string, a ...any) {
	lv := INFO
	if !l.enable(lv) {
		return
	}
	logPrint(lv, str, a...)
}

func (l *ConsoleLogger) Warn(str string, a ...any) {
	lv := WARN
	if !l.enable(lv) {
		return
	}
	logPrint(lv, str, a...)
}

func (l *ConsoleLogger) Error(str string, a ...any) {
	lv := ERROR
	if !l.enable(lv) {
		return
	}
	logPrint(lv, str, a...)
}

func (l *ConsoleLogger) Fatal(str string, a ...any) {
	lv := FATAL
	if !l.enable(lv) {
		return
	}
	logPrint(lv, str, a...)
}
