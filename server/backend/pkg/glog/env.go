package glog

import (
	"os"
	"strconv"
	"strings"
)

const (
	LOG_OUTPUT_SOURCE_ENV = "LOG_OUTPUT_SOURCE"
	LOG_OUTPUT_FORMAT_ENV = "LOG_OUTPUT_FORMAT"
	LOG_LEVEL_ENV         = "LOG_LEVEL"
	LOG_FILE_PATH_ENV     = "LOG_FILE_PATH"
)

func (a *Args) parseEnv() {
	WithOutputSource(OutputSource(parseEnvDefaultInt(LOG_OUTPUT_SOURCE_ENV, int(a.OutputSource))))(a)
	WithOutputFormat(OutputFormat(parseEnvDefaultInt(LOG_OUTPUT_FORMAT_ENV, int(a.OutputFormat))))(a)
	WithLevel(Level(parseEnvDefaultInt(LOG_LEVEL_ENV, int(a.Level))))(a)
	WithFilePath(parseEnvDefaultString(LOG_FILE_PATH_ENV, a.FilePath))
}

func parseEnvDefaultInt(env string, defaultValue int) int {
	if value := os.Getenv(env); value != "" {
		v, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue
		}
		return v
	}
	return defaultValue
}

func parseEnvDefaultString(env string, defaultValue string) string {
	if value := os.Getenv(env); value != "" {
		return strings.TrimSpace(value)
	}
	return defaultValue
}
