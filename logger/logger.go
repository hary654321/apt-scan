package logger

import (
	"fmt"
	"ias_tool_v2/config"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const (
	flag       = log.Ldate | log.Ltime | log.Lshortfile
	preDebug   = "[DEBUG]"
	preInfo    = "[INFO]"
	preWarning = "[WARNING]"
	preError   = "[ERROR]"
)

var (
	logFile        io.Writer
	debugLogger    *log.Logger
	infoLogger     *log.Logger
	warningLogger  *log.Logger
	errorLogger    *log.Logger
	defaultLogFile = filepath.Join(config.GlobalPath, "ias.log")
)

func init() {
	var (
		err     error
		logFile *os.File
	)
	if true {
		//此处需要改动改成从配置文件读取
		logFile = os.Stdout
	} else {
		logFile, err = os.OpenFile(defaultLogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
	}

	debugLogger = log.New(logFile, preDebug, flag)
	infoLogger = log.New(logFile, preInfo, flag)
	warningLogger = log.New(logFile, preWarning, flag)
	errorLogger = log.New(logFile, preError, flag)
}

func Debugf(format string, v ...interface{}) {
	debugLogger.Printf(format, v...)
}

func Infof(format string, v ...interface{}) {
	infoLogger.Printf(format, v...)
}

func Warningf(format string, v ...interface{}) {
	warningLogger.Printf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	var buf [1024]byte
	n := runtime.Stack(buf[:], true)
	fmt.Println(string(buf[:]), n)
	errorLogger.Printf(format, v...)
}
