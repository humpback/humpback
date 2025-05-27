package glog

import (
	"log/slog"
	"path/filepath"
)

type Level = slog.Level

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

type OutputSource int

const (
	OutputTypeStd  OutputSource = 0
	OutputTypeFile OutputSource = 1
)

type OutputFormat int

const (
	OutputFormatDefault OutputFormat = 0
	OutputFormatText    OutputFormat = 1
	OutputFormatJson    OutputFormat = 2
)

var (
	levelMap = map[Level]bool{
		LevelDebug: true,
		LevelInfo:  true,
		LevelWarn:  true,
		LevelError: true,
	}
	formatMap = map[OutputFormat]bool{
		OutputFormatDefault: true,
		OutputFormatText:    true,
		OutputFormatJson:    true,
	}
	outputTypeMap = map[OutputSource]bool{
		OutputTypeFile: true,
		OutputTypeStd:  true,
	}
)

var logArgs *Args

type Args struct {
	OutputSource OutputSource
	OutputFormat OutputFormat
	Level        Level
	FilePath     string
}

type Option func(*Args)

func initLogArgs(options ...Option) {
	fPath, _ := filepath.Abs("./logs/main.log")
	logArgs = &Args{
		OutputSource: OutputTypeStd,
		OutputFormat: OutputFormatDefault,
		Level:        LevelInfo,
		FilePath:     fPath,
	}
	for _, f := range options {
		f(logArgs)
	}
	logArgs.parseEnv()
}

// WithOutputSource 配置log的输出目的地是控制台还是文件，默认为控制台
func WithOutputSource(outputSource OutputSource) Option {
	return func(args *Args) {
		if _, ok := outputTypeMap[outputSource]; ok {
			args.OutputSource = outputSource
		} else {
			args.OutputSource = OutputTypeStd
		}
	}
}

// WithOutputFormat 配置log输出的格式是json还是line，默认为line
func WithOutputFormat(outputFormat OutputFormat) Option {
	return func(args *Args) {
		if _, ok := formatMap[outputFormat]; ok {
			args.OutputFormat = outputFormat
		} else {
			args.OutputFormat = OutputFormatDefault
		}
	}
}

// WithLevel 配置log输出level，默认为info
func WithLevel(level Level) Option {
	return func(args *Args) {
		if _, ok := levelMap[level]; ok {
			args.Level = level
		} else {
			args.Level = LevelInfo
		}
	}
}

// WithFilePath 配置输出文件名
func WithFilePath(path string) Option {
	return func(args *Args) {
		fpath, err := filepath.Abs(path)
		if err != nil {
			return
		}
		args.FilePath = fpath
	}
}
