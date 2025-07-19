package logger

import (
	"os"
	"runtime/debug"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var level zap.AtomicLevel = zap.NewAtomicLevelAt(zapcore.InfoLevel)

var logger *zap.SugaredLogger

func init() {
	levelEnv := os.Getenv("LOG_LEVEL")
	if len(levelEnv) > 0 {
		SetLevel(levelEnv)
	}
	encoder := getEncoder()
	stdoutCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.WarnLevel && level.Level().Enabled(lvl)
	}))
	stderrCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stderr), zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel && level.Level().Enabled(lvl)
	}))

	options := []zap.Option{}
	options = append(options, zap.ErrorOutput(zapcore.AddSync(os.Stderr))) // zap错误日志输出
	options = append(options, zap.AddCaller())                             // 增加调用者打印
	options = append(options, zap.AddStacktrace(zap.DPanicLevel))          // 打印堆栈的日志等级
	options = append(options, zap.AddCallerSkip(1))

	log := zap.New(zapcore.NewTee(stdoutCore, stderrCore), options...)
	logger = log.Sugar()
}

type LogFileConf struct {
	MaxSize    int  `json:"maxsize" yaml:"maxsize"`
	MaxAge     int  `json:"maxage" yaml:"maxage"`
	MaxBackups int  `json:"maxbackups" yaml:"maxbackups"`
	Compress   bool `json:"compress" yaml:"compress"`
}

// infoFilename 日志输出路径
// errFilname 错误日志输出路径
func SetLogFile(infoFilename string, errFilname string, conf LogFileConf) {
	if len(infoFilename) == 0 {
		panic("log file name can not empty")
	}
	infoLogger := &lumberjack.Logger{
		Filename:   infoFilename,    // 文件位置
		MaxSize:    conf.MaxSize,    // 进行切割之前,日志文件的最大大小(MB为单位)
		MaxAge:     conf.MaxAge,     // 保留旧文件的最大天数
		MaxBackups: conf.MaxBackups, // 保留旧文件的最大个数
		Compress:   conf.Compress,   // 是否压缩/归档旧文件
	}
	errLogger := &lumberjack.Logger{
		Filename:   errFilname,      // 文件位置
		MaxSize:    conf.MaxSize,    // 进行切割之前,日志文件的最大大小(MB为单位)
		MaxAge:     conf.MaxAge,     // 保留旧文件的最大天数
		MaxBackups: conf.MaxBackups, // 保留旧文件的最大个数
		Compress:   conf.Compress,   // 是否压缩/归档旧文件
	}

	encoder := getEncoder()
	var core zapcore.Core
	if infoFilename == errFilname || len(errFilname) == 0 {
		core = zapcore.NewCore(encoder, zapcore.AddSync(infoLogger), level)
	} else {
		stdoutCore := zapcore.NewCore(encoder, zapcore.AddSync(infoLogger), zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl <= zapcore.WarnLevel && level.Level().Enabled(lvl)
		}))
		stderrCore := zapcore.NewCore(encoder, zapcore.AddSync(errLogger), zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel && level.Level().Enabled(lvl)
		}))
		core = zapcore.NewTee(stdoutCore, stderrCore)
	}

	options := []zap.Option{}
	options = append(options, zap.ErrorOutput(zapcore.AddSync(errLogger))) // zap错误日志输出
	options = append(options, zap.AddCaller())                             // 增加调用者打印
	options = append(options, zap.AddStacktrace(zap.DPanicLevel))          // 打印堆栈的日志等级
	options = append(options, zap.AddCallerSkip(1))

	log := zap.New(core, options...)
	logger = log.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.999")
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// panic
// fatal
// error
// warn
// info
// debug
func SetLevel(str string) {
	l, err := zapcore.ParseLevel(str)
	if err != nil {
		Errorln("parse log level failed", err)
		return
	}
	level.SetLevel(l)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	logger.Panic(args...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	logger.Debugln(args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	logger.Infoln(args...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	logger.Warnln(args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	logger.Errorln(args...)
}

// Panicln logs a message at level Panic on the standard logger.
func Panicln(args ...interface{}) {
	logger.Panicln(args...)
}

// Fatalln logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalln(args ...interface{}) {
	logger.Fatalln(args...)
}

func DebugStackln(args ...interface{}) {
	args = append(args, "\n" + string(debug.Stack()))
	logger.Debugln(args...)
}

func InfoStackln(args ...interface{}) {
	args = append(args, "\n" + string(debug.Stack()))
	logger.Infoln(args...)
}

func WarnStackln(args ...interface{}) {
	args = append(args, "\n" + string(debug.Stack()))
	logger.Warnln(args...)
}

func ErrorStackln(args ...interface{}) {
	args = append(args, "\n" + string(debug.Stack()))
	logger.Errorln(args...)
}
