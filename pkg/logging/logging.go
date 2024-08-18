package logging

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// contextKey is a type alias for string used to store and retrieve values from context.
type contextKey = string

// loggerKey is the key used to store the logger in the context.
const loggerKey = contextKey("logger")

var (
	// defaultLogger is the default logger instance for the application.
	// It is initialized only once per package, when DefaultLogger is first called.
	defaultLogger     *zap.SugaredLogger
	defaultLoggerOnce sync.Once // Ensures defaultLogger is only initialized once.
)

// Config holds the configuration settings for the logger.
type Config struct {
	Encoding     string        // Log output format: "console" or "json"
	Level        zapcore.Level // Default log level (e.g., Info, Debug, Error)
	Development  bool          // Whether the logger is in development mode
	LogToFile    bool          // Whether to log to a file (automatically enabled in production)
	LogDirectory string        // Directory where log files will be stored
	Production   bool          // Whether the application is in production mode
}

// conf is the default logger configuration.
var conf = &Config{
	Encoding:     "console",         // Log output format: "console" or "json"
	Level:        zapcore.InfoLevel, // Default log level (Info)
	Development:  true,              // Development mode enabled by default
	LogToFile:    false,             // By default, do not log to a file
	LogDirectory: "./logs",          // Default directory for log files 	// By default, not in production mode
}

// SetConfig updates the logging configuration for the default logger.
// It automatically enables file logging if the application is in production mode.
// Must be called before DefaultLogger() to take effect.
func SetConfig(c *Config) {
	conf = &Config{
		Encoding:     c.Encoding,
		Level:        c.Level,
		Development:  c.Development,
		LogDirectory: c.LogDirectory,
	}

	// Enable file logging automatically if in production mode
	if !conf.Development {
		conf.LogToFile = true
	}
}

// SetLevel updates the logging level for the default logger.
func SetLevel(l zapcore.Level) {
	conf.Level = l
}

// NewLogger creates a new logger instance based on the provided configuration.
// It returns a SugaredLogger, which is a wrapper around zap's Logger that provides
// a more user-friendly API.
func NewLogger(conf *Config) *zap.SugaredLogger {
	// Create a default encoder configuration
	ec := zap.NewProductionEncoderConfig()
	ec.EncodeTime = zapcore.ISO8601TimeEncoder // Set time format to ISO8601

	// Set up output paths, starting with standard output
	outputPaths := []string{"stdout"}
	errorOutputPaths := []string{"stderr"}

	// If logging to a file is enabled, add a file output path
	if conf.LogToFile {
		// Ensure the log directory exists
		if err := os.MkdirAll(conf.LogDirectory, os.ModePerm); err != nil {
			fmt.Printf("Failed to create log directory: %v\n", err)
			outputPaths = append(outputPaths, "stdout") // Fallback to stdout if directory creation fails
		} else {
			// Generate the log file name based on the current date
			logFileName := filepath.Join(conf.LogDirectory, time.Now().Format("2006-01-02")+".log")
			outputPaths = append(outputPaths, logFileName)
		}
	}

	// Create the logger configuration
	cfg := zap.Config{
		Encoding:         conf.Encoding,                    // Set the log format (console or JSON)
		EncoderConfig:    ec,                               // Apply the encoder configuration
		Level:            zap.NewAtomicLevelAt(conf.Level), // Set the log level
		Development:      conf.Development,                 // Enable development mode if set
		OutputPaths:      outputPaths,                      // Log output destinations
		ErrorOutputPaths: errorOutputPaths,                 // Error log output destinations
	}

	// Build the logger and handle any errors
	logger, err := cfg.Build()
	if err != nil {
		logger = zap.NewNop() // Fallback to a no-op logger if building fails
	}
	return logger.Sugar() // Return the SugaredLogger
}

// DefaultLogger returns the default logger for the application.
// It initializes the logger only once, on the first call.
func DefaultLogger() *zap.SugaredLogger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = NewLogger(conf) // Initialize the logger with the current configuration
	})
	return defaultLogger
}

// WithLogger creates a new context with the provided logger attached to it.
// This is useful for passing the logger around in request handlers and other contexts.
func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	// If the context is a Gin context, extract the underlying request context
	if gCtx, ok := ctx.(*gin.Context); ok {
		ctx = gCtx.Request.Context()
	}
	// Attach the logger to the context
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext retrieves the logger stored in the context.
// If no logger is found in the context, the default logger is returned.
func FromContext(ctx context.Context) *zap.SugaredLogger {
	// If the context is nil, return the default logger
	if ctx == nil {
		return DefaultLogger()
	}

	// If the context is a Gin context, extract the underlying request context
	if gCtx, ok := ctx.(*gin.Context); ok && gCtx != nil {
		ctx = gCtx.Request.Context()
	}

	// Retrieve the logger from the context using the loggerKey
	if logger, ok := ctx.Value(loggerKey).(*zap.SugaredLogger); ok {
		return logger
	}

	// If no logger is found in the context, return the default logger
	return DefaultLogger()
}
