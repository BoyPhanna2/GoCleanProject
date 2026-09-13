package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"myapp/internal/config"
)

// dailyRotator implements io.WriteCloser and rotates files daily.
// It writes to logs/YYYY-MM-DD.log
type dailyRotator struct {
	mu       sync.Mutex
	dir      string
	currDate string
	file     *os.File
}

func newDailyRotator(dir string) (*dailyRotator, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	dr := &dailyRotator{dir: dir}
	if err := dr.rotate(); err != nil {
		return nil, err
	}
	return dr, nil
}

func (d *dailyRotator) rotate() error {
	now := time.Now()
	dateStr := now.Format("2006-01-02")
	if d.currDate == dateStr && d.file != nil {
		return nil
	}

	if d.file != nil {
		_ = d.file.Close()
	}

	filename := filepath.Join(d.dir, fmt.Sprintf("%s.log", dateStr))
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	d.file = f
	d.currDate = dateStr
	return nil
}

func (d *dailyRotator) Write(p []byte) (n int, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.rotate(); err != nil {
		return 0, err
	}
	return d.file.Write(p)
}

func (d *dailyRotator) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.file != nil {
		return d.file.Close()
	}
	return nil
}

func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	rotator, err := newDailyRotator(cfg.LogDir)
	if err != nil {
		return nil, err
	}

	// JSON encoder configuration
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Write to both stdout and the daily log file
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(rotator), zapcore.DebugLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	)

	return zap.New(core), nil
}

// Module provides the zap logger to Fx.
var Module = fx.Provide(
	func(cfg *config.Config) (*zap.Logger, error) {
		return NewLogger(cfg)
	},
)
