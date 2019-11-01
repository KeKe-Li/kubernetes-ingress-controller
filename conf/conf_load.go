package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/chanxuehong/log"
	"github.com/chanxuehong/log/trace"
)

func LoadConfig() *Config {
	file := getConf("conf/conf.toml")
	var configValue Config
	meta, err := toml.DecodeFile(file, &configValue)
	if err != nil {
		log.Fatal("decode config file failed", "error", err.Error())
		os.Exit(1)
	}
	mustValidate(&meta)

	// 设置 log 参数
	{
		// 设置 std logger options
		{
			if configValue.Log.Format == "json" {
				log.SetFormatter(log.JsonFormatter)
			} else {
				log.SetFormatter(log.TextFormatter)
			}
			log.SetLevelString(configValue.Log.Level)
		}

		// 设置默认的 options
		{
			defaultLogOptions := make([]log.Option, 0, 4)
			defaultLogOptions = append(defaultLogOptions, log.WithTraceIdFunc(func() string {
				return trace.NewTraceId()
			}))
			if configValue.Log.Format == "json" {
				defaultLogOptions = append(defaultLogOptions, log.WithFormatter(log.JsonFormatter))
			} else {
				defaultLogOptions = append(defaultLogOptions, log.WithFormatter(log.TextFormatter))
			}
			defaultLogOptions = append(defaultLogOptions, log.WithLevelString(configValue.Log.Level))
			log.SetDefaultOptions(defaultLogOptions)
		}
	}

	return &configValue
}

func getConf(file string) (full string) {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
	for {
		full = filepath.Join(path, file)
		if info, err := os.Stat(full); err != nil {
			// other erros
			if !os.IsNotExist(err) {
				log.Fatal(err.Error(), "full", full)
				os.Exit(1)
			}
		} else if info != nil && info.IsDir() {
			// exit but is directory
		} else {
			return full
		}
		i := strings.LastIndexByte(path, filepath.Separator)
		if i == -1 {
			log.Fatal("file not exist")
			os.Exit(1)
		}
		path = path[:i]
	}
}

func mustValidate(meta *toml.MetaData) {
	for _, v := range requires {
		if !meta.IsDefined(strings.Split(v, ".")...) {
			log.Fatal(fmt.Sprintf("%s required", v))
			os.Exit(1)
		}
	}
}
