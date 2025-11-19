package logx

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"

	"github.com/suhanyujie/go_utils/env_x"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 自定义日志库 todo
var (
	Logger      *Logx
	defaultPath = "logs/test.log"
)

type LogConfIf interface {
	GetLogDir() string
}

type Logx struct {
	inner *zap.Logger
}

func NewLogger() *Logx {
	l, err := zap.NewDevelopment()
	if err != nil {
		log.Printf("[NewLogger] err: %v", err)
		panic(err)
	}
	ins := &Logx{
		inner: l,
	}
	return ins
}

func InitLogger(l *Logx) {
	if Logger == nil {
		Logger = l
	}
}

func GetLogger() *Logx {
	return Logger
}

func GetSysLogger(confObj LogConfIf) *Logx {
	if Logger == nil {
		once := sync.Once{}
		once.Do(func() {
			Logger = NewLoggerForSys(confObj.GetLogDir())
		})
	}
	return Logger
}

func NewLoggerForRoom() *zap.Logger {
	l, err := zap.NewDevelopment()
	if err != nil {
		log.Printf("[NewLogger] err: %v", err)
		panic(err)
	}

	return l
}

// 定制实例化
func NewFromCustom() *Logx {
	encoderConfig := zap.NewProductionEncoderConfig()
	// 设置日志记录中时间的格式
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	// 日志 Encoder 还是 JSONEncoder，把日志行格式化成JSON格式的
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	fileWriteSyncer := getFileLogWriter(defaultPath)

	core := zapcore.NewTee(
		// 同时向控制台和文件写日志， 生产环境记得把控制台写入去掉，日志记录的基本是 Debug 及以上，生产环境记得改成 Info
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(encoder, fileWriteSyncer, zapcore.DebugLevel),
	)
	logger := zap.New(core)
	logger = logger.WithOptions(zap.Hooks(FsNoticeHook()))
	ins := &Logx{
		inner: logger,
	}

	return ins
}

// 实例化某个游戏房间的日志记录器
func NewByRoom(rid uint64, dir string) *Logx {
	if dir == "" {
		dir = "./log/game/"
	}
	logFilePah := dir
	encoderConfig := zap.NewProductionEncoderConfig()

	// 设置日志记录中时间的格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 日志 Encoder 还是 JSONEncoder，把日志行格式化成JSON格式的
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	encoder.AddUint64("rid", rid)
	fileWriteSyncer := getFileLogWriter(logFilePah)

	core := zapcore.NewTee(
		// 同时向控制台和文件写日志， 生产环境记得把控制台写入去掉，日志记录的基本是 Debug 及以上，生产环境记得改成 Info
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(encoder, fileWriteSyncer, zapcore.DebugLevel),
	)
	logger := zap.New(core)
	logger = logger.WithOptions(zap.Hooks(FsNoticeHook()))
	ins := &Logx{
		inner: logger,
	}

	return ins
}

func NewLoggerForSys(dir string) *Logx {
	if dir == "" {
		dir = "./log/game/app.log"
	}
	logFilePah := dir
	encoderConfig := zap.NewProductionEncoderConfig()

	// 设置日志记录中时间的格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 日志 Encoder 还是 JSONEncoder，把日志行格式化成JSON格式的
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	encoder.AddString("from", "sys")
	fileWriteSyncer := getFileLogWriter(logFilePah)

	core := zapcore.NewTee(
		// 同时向控制台和文件写日志， 生产环境记得把控制台写入去掉，日志记录的基本是 Debug 及以上，生产环境记得改成 Info
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(encoder, fileWriteSyncer, zapcore.DebugLevel),
	)
	logger := zap.New(core)
	logger = logger.WithOptions(zap.Hooks(FsNoticeHook()))
	ins := &Logx{
		inner: logger,
	}

	return ins
}

func NewLoggerByTag(dir string, tag string) *Logx {
	if dir == "" {
		dir = "./log/game/app.log"
	}
	logFilePah := dir
	encoderConfig := zap.NewProductionEncoderConfig()

	// 设置日志记录中时间的格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 日志 Encoder 还是 JSONEncoder，把日志行格式化成JSON格式的
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	encoder.AddString("from", tag)
	fileWriteSyncer := getFileLogWriter(logFilePah)

	core := zapcore.NewTee(
		// 同时向控制台和文件写日志， 生产环境记得把控制台写入去掉，日志记录的基本是 Debug 及以上，生产环境记得改成 Info
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(encoder, fileWriteSyncer, zapcore.DebugLevel),
	)
	logger := zap.New(core)
	logger = logger.WithOptions(zap.Hooks(FsNoticeHook()))
	ins := &Logx{
		inner: logger,
	}

	return ins
}

func getFileLogWriter(logPath string) (writeSyncer zapcore.WriteSyncer) {
	if logPath == "" {
		logPath = "./log/game/app.log"
	}
	// 使用 lumberjack 实现 log rotate
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    100, // 单个文件最大 100M
		MaxBackups: 10,  // 多于 n 个日志文件后，清理较旧的日志
		MaxAge:     15,  // 保留的天数
		Compress:   false,
	}

	return zapcore.AddSync(lumberJackLogger)
}

// ref https://juejin.cn/post/7110892495827009572
// n 一般为 2
func getCallerInfoForLog(n int) (callerFields []zap.Field) {
	pc, file, line, ok := runtime.Caller(n) // 回溯 2 层，拿到写日志的调用方的函数信息
	if !ok {
		return
	}
	funcName := runtime.FuncForPC(pc).Name()
	funcName = path.Base(funcName) // Base函数返回路径的最后一个元素，只保留函数名

	callerFields = append(callerFields, zap.String("func", funcName), zap.String("file", file), zap.Int("line", line))
	return
}

func FsNoticeHook() func(zapcore.Entry) error {
	return func(entry zapcore.Entry) error {
		if entry.Level < zapcore.ErrorLevel {
			return nil
		}
		// 临时方式。后续最好集成到 sentry 中。
		// 发送告警到飞书群 todo
		env := env_x.GetEnv(env_x.DefaultEnvKey)
		_, err := Send2FsWithWarningMsg(env, entry.Message)
		return err
	}
}

func SimplePost(url string, inputBody string) (string, int, error) {
	resp, err := http.Post(url,
		"application/json",
		strings.NewReader(inputBody))
	if err != nil {
		Logger.Infof("[SimplePost] Post err:%v", err)
		return "", -201, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logger.Infof("[SimplePost] ReadAll err:%v", err)
		return "", resp.StatusCode, err
	}
	bodyStr := string(respBody)
	return bodyStr, resp.StatusCode, nil
}

/// 以下是打印日志的一些方法

func (l *Logx) Info(msg string, fields ...zap.Field) {
	l.inner.Info(msg, fields...)
}

func (l *Logx) Error(msg string, fields ...zap.Field) {
	l.inner.Error(msg, fields...)
}

func (l *Logx) Errorf(fmtStr string, params ...any) {
	msg := fmt.Sprintf(fmtStr, params...)
	callerInfoArr := getCallerInfoForLog(2)
	l.inner.Error(msg, callerInfoArr...)
}

func (l *Logx) Infof(fmtStr string, params ...any) {
	msg := fmt.Sprintf(fmtStr, params...)
	callerInfoArr := getCallerInfoForLog(2)
	l.inner.Info(msg, callerInfoArr...)
}

func (l *Logx) infoFInner(fmtStr string, params ...any) {
	msg := fmt.Sprintf(fmtStr, params...)
	// 被包内的其他方法调用，因此 skip 层数为 3
	callerInfoArr := getCallerInfoForLog(3)
	l.inner.Info(msg, callerInfoArr...)
}

// 推荐使用 Infof / Errorf / Debugf
func (l *Logx) Debug(msg string, fields ...zap.Field) {
	l.inner.Debug(msg, fields...)
}

func (l *Logx) Debugf(fmtStr string, params ...any) {
	msg := fmt.Sprintf(fmtStr, params...)
	callerInfoArr := getCallerInfoForLog(2)
	l.inner.Debug(msg, callerInfoArr...)
}

func (l *Logx) debugFInner(fmtStr string, params ...any) {
	msg := fmt.Sprintf(fmtStr, params...)
	// 被包内的其他方法调用，因此 skip 层数为 3
	callerInfoArr := getCallerInfoForLog(3)
	l.inner.Debug(msg, callerInfoArr...)
}

// 推荐使用 Infof / Errorf
func (l *Logx) Printf(fmtStr string, params ...any) {
	// l.Infof(fmtStr, params...)
	l.infoFInner(fmtStr, params...)
}

func (l *Logx) Println(params ...any) {
	str := strings.Builder{}
	for _, item := range params {
		str.WriteString(fmt.Sprintf("%v", item))
	}
	msg := str.String()
	callerInfoArr := getCallerInfoForLog(2)
	l.inner.Info(msg, callerInfoArr...)
}

func (l *Logx) Fatal(params ...any) {
	str := strings.Builder{}
	for _, item := range params {
		str.WriteString(fmt.Sprintf("%v", item))
	}
	msg := str.String()
	callerInfoArr := getCallerInfoForLog(2)
	l.inner.Info(msg, callerInfoArr...)
}

func (l *Logx) Fatalf(fmtStr string, params ...any) {
	l.Errorf(fmtStr, params...)
}

func (l *Logx) GetInnerLogger() *zap.Logger {
	return l.inner
}
