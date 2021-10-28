package main

import (
	log "github.com/imb2022/log4go"
)

// SetLog set logger
func SetLog() {
	w := log.NewConsoleWriterWithLevel(log.DEBUG)
	w.SetColor(true)

	log.Register(w)
	log.SetLayout("2006-01-02 15:04:05")
}

func main() {
	SetLog()
	defer log.Close()

	var name = "xwi88"
	log.Debug("log4go by %s", name)
	log.Info("log4go by %s", name)
	log.Warn("log4go by %s", name)
	log.Error("log4go by %s", name)
	log.Fatal("log4go by %s", name)
}
