package kubesec

import (
	"errors"
	"fmt"
	"io"

	"github.com/controlplaneio/kubesec/v2/pkg/ruler"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Client represent a client for kubesec.
type Client struct {
}

// NewClient returns a new client for kubesec.
func NewClient() *Client {
	return &Client{}
}

func newLogger(logLevel string, zapEncoding string) (*zap.SugaredLogger, error) {
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	switch logLevel {
	case "debug":
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	case "panic":
		level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	}

	zapEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "severity",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	zapConfig := zap.Config{
		Level:       level,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         zapEncoding,
		EncoderConfig:    zapEncoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}

// ScanDefinition scans the provided resource definition.
func (kc *Client) ScanDefinition(def []byte) (*ruler.Report, error) {
	var logger *zap.SugaredLogger
	logger, err := newLogger("info", "console")
	if err != nil {
		return nil, err
	}

	results, err := ruler.NewRuleset(logger).Run("SCAN", def, "")
	if err != nil {
		return nil, err
	}

	if len(results) != 1 {
		return nil, errors.New("Unexpected amount of results")
	}
	result := results[0]

	return &result, nil
}

// DumpReport writes the result in a human-readable format to the specified writer.
func DumpReport(r *ruler.Report, w io.Writer) {
	fmt.Fprintf(w, "kubesec score: %v\n", r.Score)
	fmt.Fprintln(w, "-----------------")
	if len(r.Scoring.Critical) > 0 {
		fmt.Fprintln(w, "Critical")
		for i, el := range r.Scoring.Critical {
			fmt.Fprintf(w, "%v. %v\n", i+1, el.Selector)
			if len(el.Reason) > 0 {
				fmt.Fprintln(w, el.Reason)
			}

		}
		fmt.Fprintln(w, "-----------------")
	}
	if len(r.Scoring.Advise) > 0 {
		fmt.Fprintf(w, "Advise")
		for i, el := range r.Scoring.Advise {
			fmt.Fprintf(w, "%v. %v\n", i+1, el.Selector)
			if len(el.Reason) > 0 {
				fmt.Fprintln(w, el.Reason)
			}
		}
	}
}
