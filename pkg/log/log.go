package log

import (
	"io"
	"sync"
)

var (
	loggerConfig Config

	_debugLogger  *Logger
	_debugLogOnce sync.Once
	_infoLogger   *Logger
	_infoLogOnce  sync.Once
	_warnLogger   *Logger
	_warnLogOnce  sync.Once
	_errorLogger  *Logger
	_errorLogOnce sync.Once
	_traceLogger  *Logger
	_traceLogOnce sync.Once
)

func debugLogger() *Logger {
	_debugLogOnce.Do(func() {
		_debugLogger = New(loggerConfig, "debug")
	})
	return _debugLogger
}

func infoLogger() *Logger {
	_infoLogOnce.Do(func() {
		_infoLogger = New(loggerConfig, "info")
	})
	return _infoLogger
}

func warnLogger() *Logger {
	_warnLogOnce.Do(func() {
		_warnLogger = New(loggerConfig, "warn")
	})
	return _warnLogger
}

func errorLogger() *Logger {
	_errorLogOnce.Do(func() {
		_errorLogger = New(loggerConfig, "error")
	})
	return _errorLogger
}

func traceLogger() *Logger {
	_traceLogOnce.Do(func() {
		_traceLogger = New(loggerConfig, "trace")
	})
	return _traceLogger
}

func InitLogger(c Config) {
	c.CallerSkip = 2
	loggerConfig = c
}

func Debug(args ...interface{}) {
	debugLogger().Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	debugLogger().Debugf(format, args...)
}

func Info(args ...interface{}) {
	infoLogger().Info(args...)
}

func Infof(format string, args ...interface{}) {
	infoLogger().Infof(format, args...)
}

func Warn(args ...interface{}) {
	warnLogger().Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	warnLogger().Warnf(format, args...)
}

func Error(args ...interface{}) {
	errorLogger().Error(args...)
}

func Errorf(format string, args ...interface{}) {
	errorLogger().Errorf(format, args...)
}

func Trace(requestID string, tag string, event string, msg string) {
	traceLogger().Trace(requestID, tag, event, msg)
}

func Tracef(requestID string, tag string, event string, format string, args ...interface{}) {
	traceLogger().Tracef(requestID, tag, event, format, args...)
}

// GetAccessWriter 获取Access输出流
func GetAccessWriter() io.Writer {
	// return _logger.w
	loggerConfig.defaultChange()
	fn := loggerConfig.LogPath("access")
	return getLogWriter(fn, loggerConfig.MaxSize)
}
