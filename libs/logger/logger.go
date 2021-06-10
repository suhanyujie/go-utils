package logger

/// 深度参考：github.com/nacos-group/nacos-sdk-go/common/logger
import (
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/nacos-group/nacos-sdk-go/common/file"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger  Logger
	logLock sync.RWMutex
)

var levelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
}

type Config struct {
	Level        string
	OutputPath   string
	RotationTime string
	MaxAge       int64
}

type MyLogger struct {
	Logger
}

// Logger is the interface for Logger types
type Logger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})

	Infof(fmt string, args ...interface{})
	Warnf(fmt string, args ...interface{})
	Errorf(fmt string, args ...interface{})
	Debugf(fmt string, args ...interface{})
}

func init() {
	zapLoggerConfig := zap.NewDevelopmentConfig()
	zapLoggerEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	zapLoggerConfig.EncoderConfig = zapLoggerEncoderConfig
	zapLogger, _ := zapLoggerConfig.Build(zap.AddCallerSkip(1))
	SetLogger(&MyLogger{zapLogger.Sugar()})
}

// InitLogger is init global logger for my service
func InitLogger(config Config) (err error) {
	logLock.Lock()
	defer logLock.Unlock()
	logger, err = InitMyLogger(config)
	return
}

// InitMyLogger is init default logger
func InitMyLogger(config Config) (Logger, error) {
	logLevel := getLogLevel(config.Level)
	encoder := getEncoder()
	writer, err := getWriter(config.OutputPath, config.RotationTime, config.MaxAge)
	if err != nil {
		return nil, err
	}
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoder), zapcore.AddSync(writer), logLevel)
	zaplogger := zap.New(core, zap.AddCallerSkip(1))
	return &MyLogger{zaplogger.Sugar()}, nil
}

func getLogLevel(level string) zapcore.Level {
	if zapLevel, ok := levelMap[level]; ok {
		return zapLevel
	}
	return zapcore.InfoLevel
}

func getEncoder() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func getWriter(outputPath string, rotateTime string, maxAge int64) (writer io.Writer, err error) {
	err = file.MkdirIfNecessary(outputPath)
	if err != nil {
		return
	}
	outputPath = outputPath + string(os.PathSeparator)
	rotateDuration, err := time.ParseDuration(rotateTime)
	writer, err = rotatelogs.New(filepath.Join(outputPath, "goPerSvc.log-%Y%m%d%H%M"),
		rotatelogs.WithRotationTime(rotateDuration), rotatelogs.WithMaxAge(time.Duration(maxAge)*rotateDuration),
		rotatelogs.WithLinkName(filepath.Join(outputPath, "goPerSvc.log")))
	return
}

//SetLogger sets logger for sdk
func SetLogger(log Logger) {
	logLock.Lock()
	defer logLock.Unlock()
	logger = log
}

func GetLogger() Logger {
	logLock.RLock()
	defer logLock.RUnlock()
	return logger
}