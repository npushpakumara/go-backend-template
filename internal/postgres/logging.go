package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/npushpakumara/go-backend-template/pkg/logging"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// msgPrefix is a prefix added to all log messages to indicate they are database-related logs.
const msgPrefix = "[DB] "

// Logger is a custom logger that implements GORM's logging interface.
// It wraps a zap.SugaredLogger for structured logging.
type Logger struct {
	cfg glogger.Config // Configuration for the logger, including log levels and thresholds.
}

// NewLogger creates and returns a new Logger instance for GORM.
// It takes the slow SQL threshold, whether to ignore "record not found" errors, and the log level as inputs.
func NewLogger(slowThreshold time.Duration, ignoreRecordNotFoundError bool, level zapcore.Level) *Logger {
	// Set up the logger configuration.
	cfg := glogger.Config{
		SlowThreshold:             slowThreshold,             // Threshold for slow SQL logging.
		Colorful:                  false,                     // Disable colorful logs (not supported by zap).
		IgnoreRecordNotFoundError: ignoreRecordNotFoundError, // Whether to ignore "record not found" errors in logs.
	}

	// Set the appropriate log level based on the zapcore level.
	switch level {
	case zapcore.DebugLevel, zapcore.InfoLevel:
		cfg.LogLevel = glogger.Info
	case zapcore.WarnLevel:
		cfg.LogLevel = glogger.Warn
	case zapcore.ErrorLevel:
		cfg.LogLevel = glogger.Error
	default:
		cfg.LogLevel = glogger.Silent
	}

	// Return the new Logger instance.
	return &Logger{cfg: cfg}
}

// LogMode sets the log level for the logger and returns a new logger instance with this configuration.
func (l *Logger) LogMode(level glogger.LogLevel) glogger.Interface {
	newLogger := *l                // Create a copy of the current logger.
	newLogger.cfg.LogLevel = level // Set the log level in the copy.
	return &newLogger              // Return the new logger with the updated log level.
}

// Info logs informational messages. It only logs if the log level is set to Info or higher.
func (l *Logger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.cfg.LogLevel >= glogger.Info {
		l.fromContext(ctx).Infof(msgPrefix+msg, args...)
	}
}

// Warn logs warning messages. It only logs if the log level is set to Warn or higher.
func (l *Logger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if l.cfg.LogLevel >= glogger.Warn {
		l.fromContext(ctx).Warnf(msgPrefix+msg, args...)
	}
}

// Error logs error messages. It only logs if the log level is set to Error or higher.
func (l *Logger) Error(ctx context.Context, msg string, args ...interface{}) {
	if l.cfg.LogLevel >= glogger.Error {
		l.fromContext(ctx).Errorf(msgPrefix+msg, args...)
	}
}

// Trace logs SQL queries and their execution times, as well as any errors that occurred.
// It is used by GORM to log the details of each SQL operation.
func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// If the logger is set to silent, do nothing.
	if l.cfg.LogLevel == glogger.Silent {
		return
	}

	elapsed := time.Since(begin) // Calculate the time taken for the SQL query.
	logger := l.fromContext(ctx) // Get the logger from the context.

	// Log formats for different scenarios.
	const (
		traceStr     = msgPrefix + "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = msgPrefix + "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = msgPrefix + "%s %s\n[%.3fms] [rows:%v] %s"
	)

	// Log the SQL query and its details based on the log level and whether an error occurred.
	switch {
	case err != nil && l.cfg.LogLevel >= glogger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.cfg.IgnoreRecordNotFoundError):
		// Log errors if any, except for "record not found" errors if configured to ignore them.
		sql, rows := fc()
		if rows == -1 {
			logger.Errorf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logger.Errorf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.cfg.SlowThreshold && l.cfg.SlowThreshold != 0 && l.cfg.LogLevel >= glogger.Warn:
		// Log slow SQL queries if they exceed the configured slow threshold.
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.cfg.SlowThreshold)
		if rows == -1 {
			logger.Warnf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logger.Warnf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.cfg.LogLevel == glogger.Info:
		// Log general SQL query information.
		sql, rows := fc()
		if rows == -1 {
			logger.Infof(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logger.Infof(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

// fromContext retrieves a zap.SugaredLogger from the provided context.
// This allows the logger to be used in a context-aware way.
func (l *Logger) fromContext(ctx context.Context) *zap.SugaredLogger {
	// Get the logger from the context and adjust the caller skip to account for wrapping.
	return logging.FromContext(ctx).WithOptions(zap.AddCallerSkip(3))
}
