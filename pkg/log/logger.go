package log

import (
	"fmt"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 内部实现使用zap.SugaredLogger
// 实现apiserver.Logger接口
// 日志级别 Error > Warn > Info > Debug

type (
	Logger struct {
		sugar     *zap.SugaredLogger
		w         zapcore.WriteSyncer
		serviceID string
	}

	Tracer struct {
		ReqeustID string `json:"request_id"`
		Event     string `json:"event"`
		Tag       string `json:"tag"`
		Msg       string `json:"msg"`
	}
)

func (l *Logger) Debug(args ...interface{}) {
	l.sugar.Debug(args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.sugar.Debugf(format, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.sugar.Info(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.sugar.Infof(format, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.sugar.Warn(args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.sugar.Warnf(format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.sugar.Error(args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.sugar.Errorf(format, args...)
}

// Trace trace过滤所有日志
func (l *Logger) Trace(requestID string, tag string, event string, msg string) {
	// l.w.Write()
	trageInfo := Tracer{
		ReqeustID: requestID,
		Tag:       tag,
		Event:     event,
		Msg:       msg,
	}

	fmt.Fprintf(l.w, "{\"level\":\"TRACE\", \"ts\":\"%s\", \"request_id\":\"%s\", \"event\":\"%s\", \"tag\":\"%s\", \"msg\":\"%s\", \"service_id\":\"%s\"}\n", time.Now().Format("2006-01-02T15:04:05.000Z"), trageInfo.ReqeustID, trageInfo.Event, trageInfo.Tag, trageInfo.Msg, l.serviceID)
}

func (l *Logger) Tracef(requestID string, tag string, event string, format string, args ...interface{}) {
	l.Trace(requestID, tag, event, fmt.Sprintf(format, args...))
}

// New 通过自定义输出流，限制文件输出大小控制
// 未设置自动清理
// 未实现按日分文件
// 使用LocalTime作为分文件的时间点
// 旧文件格式 /var/log/foo/server-2016-11-04T18-30-00.000.log
func New(c Config) *Logger {
	c.defaultChange()
	fn := c.LogPath()
	writeSyncer := getLogWriter(fn, c.MaxSize)

	encoder := getEncoder()

	core := zapcore.NewCore(encoder, writeSyncer, getLevel(c.Level))
	_logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(c.CallerSkip), zap.AddStacktrace(zap.WarnLevel), zap.Fields(zap.Field(zap.String("service_id", c.ServiceID()))))

	zap.RedirectStdLog(_logger)

	return &Logger{sugar: _logger.Sugar(), w: writeSyncer, serviceID: c.ServiceID()}
}

// getLevel 转换级别，将关键字转换为zap的常量
func getLevel(level string) zapcore.Level {
	var _zapLevel zapcore.Level
	switch level {
	case "debug":
		_zapLevel = zapcore.DebugLevel
	case "info":
		_zapLevel = zapcore.InfoLevel
	case "warn":
		_zapLevel = zapcore.WarnLevel
	case "error":
		_zapLevel = zapcore.ErrorLevel
	}
	return _zapLevel
}

// getLogWriter 设置保存日志文件的条件
func getLogWriter(fn string, maxsize int) zapcore.WriteSyncer {
	// MaxBackups: 3,  // 最多保存日志文件数量
	// MaxAge:    1,   // 保存旧日志最大天数
	lumberJackLogger := &lumberjack.Logger{
		Filename:  fn,
		MaxSize:   maxsize,
		LocalTime: true,
		Compress:  false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	return zapcore.NewJSONEncoder(encoderConfig)
}
