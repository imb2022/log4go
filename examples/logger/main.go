package main

import (
	"fmt"

	"github.com/imb2022/log4go"
)

func SetLog() (logger *log4go.Logger) {
	w := log4go.NewFileWriter()
	if err := w.SetPathPattern("/tmp/logs/error%Y%M%D%H.log"); err != nil {
		fmt.Printf("file writer SetPathPattern err:%v\n", err)
		return
	}

	logger = log4go.NewLogger()
	logger.Register(w)
	logger.SetLevel(log4go.DEBUG) // global min level
	return
}

func main() {
	logger := SetLog()
	defer logger.Close()

	var name = "xwi88"
	logger.Debug("log4go by %s", name)
	logger.Info("log4go by %s", name)
	logger.Warn("log4go by %s", name)
	logger.Error("log4go by %s", name)
	logger.Fatal("log4go by %s", name)
}
