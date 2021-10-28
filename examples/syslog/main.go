package main

import (
	log "github.com/imb2022/log4go"
)

func SetLog() {
	w := log.NewSyslogWriter()
	w.SetNetwork("udp")
	w.SetAddr("127.0.0.1:514")
	w.SetTag("log4go")

	log.Register(w)
	log.SetLevel(log.DEBUG)
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
