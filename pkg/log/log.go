package log

import "io"

var _logger *Logger

func InitLogger(c Config) {
	c.CallerSkip = 2
	_logger = New(c)
}

func Debug(args ...interface{}) {
	_logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	_logger.Debugf(format, args...)
}

func Info(args ...interface{}) {
	_logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	_logger.Infof(format, args...)
}

func Warn(args ...interface{}) {
	_logger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	_logger.Warnf(format, args...)
}

func Error(args ...interface{}) {
	_logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	_logger.Errorf(format, args...)
}

// GetWriter 获取输出流
func GetWriter() io.Writer {
	return _logger.w
}
