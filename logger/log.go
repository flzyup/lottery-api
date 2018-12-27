package logger

import (
	"github.com/op/go-logging"
	"os"
)

var loggers map[string]*logging.Logger

var format = logging.MustStringFormatter(
	`%{color}%{time:2006-01-02 15:04:05.000} %{shortfunc} ▶ %{level:.4s}%{color:reset} %{message}`,
	//`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

/*
 * 获得对应的logger对象
 *
 * @param name
 */
func GetLogger(name string) *logging.Logger {
	if v, ok := loggers[name]; ok {
		return v
	} else {
		logger := logging.MustGetLogger(name)
		loggers[name] = logger
		return logger
	}
}

/*
 * 初始化方法
 */
func init() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)

	formatter := logging.NewBackendFormatter(backend, format)

	logging.SetBackend(formatter)

	loggers = make(map[string]*logging.Logger)
}
