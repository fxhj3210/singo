package util

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"

	"os"
	"time"
)

const (
	// LevelError 错误
	LevelError = iota
	// LevelWarning 警告
	LevelWarning
	// LevelInformational 提示
	LevelInformational
	// LevelDebug 除错
	LevelDebug
)

var logger *Loggers

// Loggers 日志
type Loggers struct {
	level int
}

// Println 打印
func (ll *Loggers) Println(msg string) {
	fmt.Printf("%s %s\n", time.Now().Format("2006-01-02 15:04:05 -0700"), msg)
}

// Panic 极端错误
func (ll *Loggers) Panic(format string, v ...interface{}) {
	if LevelError > ll.level {
		return
	}
	msg := fmt.Sprintf("[Panic] "+format, v...)
	log.Panic(msg)
	ll.Println(msg)
	os.Exit(0)
}

// Error 错误
func (ll *Loggers) Error(format string, v ...interface{}) {
	if LevelError > ll.level {
		return
	}
	msg := fmt.Sprintf("[E] "+format, v...)
	log.Error(msg)
	ll.Println(msg)
}

// Warning 警告
func (ll *Loggers) Warning(format string, v ...interface{}) {
	if LevelWarning > ll.level {
		return
	}
	msg := fmt.Sprintf("[W] "+format, v...)
	log.Warning(msg)
	ll.Println(msg)
}

// Info 信息
func (ll *Loggers) Info(format string, v ...interface{}) {
	if LevelInformational > ll.level {
		return
	}
	msg := fmt.Sprintf("[I] "+format, v...)
	log.Info(msg)
	ll.Println(msg)
}

// Debug 校验
func (ll *Loggers) Debug(format string, v ...interface{}) {
	if LevelDebug > ll.level {
		return
	}
	msg := fmt.Sprintf("[D] "+format, v...)
	log.Debug(msg)
	ll.Println(msg)
}

// BuildLogger 构建logger
func BuildLogger(level string) {
	intLevel := LevelError
	switch level {
	case "error":
		intLevel = LevelError
	case "warning":
		intLevel = LevelWarning
	case "info":
		intLevel = LevelInformational
	case "debug":
		intLevel = LevelDebug
	default:
		intLevel = LevelDebug
	}
	l := Loggers{
		level: intLevel,
	}
	logger = &l
	path := "runtime/go.log"

	/* 日志轮转相关函数
	`WithLinkName` 为最新的日志建立软连接
	`WithRotationTime` 设置日志分割的时间，隔多久分割一次
	`WithMaxAge 和 WithRotationCount二者只能设置一个
	  `WithMaxAge` 设置文件清理前的最长保存时间
	  `WithRotationCount` 设置文件清理前最多保存的个数
	*/
	writer, _ := rotatelogs.New(
		path+".%Y%m%d",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithRotationCount(3),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	log.SetOutput(writer)
	log.SetFormatter(&log.JSONFormatter{})
}

// Log 返回日志对象
func Log() *Loggers {
	if logger == nil {
		l := Loggers{
			level: LevelDebug,
		}
		logger = &l
	}
	return logger
}
