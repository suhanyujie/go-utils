package env_x

import (
	"os"
	"sync"
)

var (
	RunENV        string
	DefaultEnvKey = "RUN_ENV"
)

func GetEnv(envKey string) string {
	if envKey == "" {
		envKey = DefaultEnvKey
	}
	if RunENV == "" {
		once := sync.Once{}
		once.Do(func() {
			runEnv := os.Getenv(envKey)
			if runEnv == "" {
				runEnv = "dev"
			}
			RunENV = runEnv
		})
	}

	return RunENV
}

// IsDev 检查是否是开发环境
func IsDev() bool {
	envVal := GetEnv(DefaultEnvKey)
	return envVal == "dev"
}
