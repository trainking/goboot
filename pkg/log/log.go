package log

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
	_logger.sugar.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	_logger.sugar.Warnf(format, args...)
}

func Error(args ...interface{}) {
	_logger.sugar.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	_logger.sugar.Errorf(format, args...)
}
