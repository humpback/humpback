package glog

import (
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

var (
	file       *os.File
	replaceStd = func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
		}
		return a
	}

	replaceFile = func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			_, file, line, ok := runtime.Caller(6)
			if ok {
				source.File = filepath.Base(file)
				source.Line = line
			}
		}
		return a
	}
)

func Open(options ...Option) error {
	initLogArgs(options...)
	return setResource()
}

func Close() {
	if file != nil {
		file.Close()
		file = nil
	}
}

func setResource() error {
	var err error
	if logArgs.OutputSource == OutputTypeFile {
		file, err = openFile(logArgs.FilePath)
		if err != nil {
			return err
		}
		initLog(file)
		return nil
	}
	initLog(os.Stderr)
	return nil
}

func initLog(w io.Writer) {
	var replaceFunc = replaceStd
	if logArgs.OutputSource == OutputTypeFile {
		replaceFunc = replaceFile
	}
	switch logArgs.OutputFormat {
	case OutputFormatText:
		slog.SetDefault(slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{
			AddSource:   true,
			Level:       logArgs.Level,
			ReplaceAttr: replaceFunc,
		})))
	case OutputFormatJson:
		slog.SetDefault(slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
			AddSource:   true,
			Level:       logArgs.Level,
			ReplaceAttr: replaceFunc,
		})))
	default:
		log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
		log.SetOutput(w)
	}
}

func openFile(filePath string) (*os.File, error) {
	if err := mkdir(filePath); err != nil {
		return nil, err
	}
	return os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
}

func mkdir(path string) error {
	dirPath := filepath.Dir(path)
	if pathExist(dirPath) {
		return nil
	}
	return os.Mkdir(dirPath, 0755)
}

func pathExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return os.IsExist(err)
	}
	return true
}
