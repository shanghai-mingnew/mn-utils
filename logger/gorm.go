package logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	gLog "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type gormLogger struct {
	SlowThreshold time.Duration

	log *zap.SugaredLogger
}

// LogMode log mode
func (l *gormLogger) LogMode(level gLog.LogLevel) gLog.Interface {
	return l
}

// Info print info
func (l *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.log.Debugln(utils.FileWithLineNum(), msg, data)
}

// Warn print warn messages
func (l *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.log.Warnln(utils.FileWithLineNum(), msg, data)
}

// Error print error messages
func (l *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.log.Errorln(utils.FileWithLineNum(), msg, data)
}

// Trace print sql message
//
//nolint:cyclop
func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case err != nil && level.Level().Enabled(zapcore.ErrorLevel) && (!errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		if rows == -1 {
			l.log.Errorf("%s %s [%.3fms] [rows:%v] \n%s", utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.log.Errorf("%s %s [%.3fms] [rows:%v] \n%s", utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && level.Level().Enabled(zapcore.WarnLevel):
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.log.Warnf("%s %s [%.3fms] [rows:%v] \n%s", utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.log.Warnf("%s %s [%.3fms] [rows:%v] \n%s", utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case level.Level().Enabled(zapcore.DebugLevel):
		sql, rows := fc()
		if rows == -1 {
			l.log.Debugf("%s [%.3fms] [rows:%v] \n%s", utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.log.Debugf("%s [%.3fms] [rows:%v] \n%s", utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

func NewGormLog(SlowThreshold time.Duration) *gormLogger {
	return &gormLogger{
		SlowThreshold: SlowThreshold,
		log:           logger.WithOptions(zap.WithCaller(false)),
	}
}
