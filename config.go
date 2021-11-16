package log4go

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"github.com/imb2022/log4go/util"
)

// GlobalLevel global level
var GlobalLevel = DEBUG

const (
	WriterNameConsole = "writer_console"
	WriterNameFile    = "writer_file"
	WriterNameKafka   = "writer_kafka"
)

// ConfFileWriter file writer config
type ConfFileWriter struct {
	Level   string `json:"level" mapstructure:"level"`
	LogPath string `json:"log_path" mapstructure:"log_path"`
	Enable  bool   `json:"enable" mapstructure:"enable"`
}

// ConfConsoleWriter console writer config
type ConfConsoleWriter struct {
	Level  string `json:"level" mapstructure:"level"`
	Enable bool   `json:"enable" mapstructure:"enable"`
	Color  bool   `json:"color" mapstructure:"color"`
}

// ConfAliLogHubWriter ali log hub writer config
type ConfAliLogHubWriter struct {
	Level           string `json:"level" mapstructure:"level"`
	Enable          bool   `json:"enable" mapstructure:"enable"`
	LogName         string `json:"log_name" mapstructure:"log_name"`
	LogSource       string `json:"log_source" mapstructure:"log_source"`
	ProjectName     string `json:"project_name" mapstructure:"project_name"`
	Endpoint        string `json:"endpoint" mapstructure:"endpoint"`
	AccessKeyId     string `json:"access_key_id" mapstructure:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret" mapstructure:"access_key_secret"`
	StoreName       string `json:"store_name" mapstructure:"store_name"`
	BufSize         int    `json:"buf_size" mapstructure:"buf_size"`
}

// LogConfig log config
type LogConfig struct {
	// global level, maybe override by real minimum multi writer level
	Level         string            `json:"level" mapstructure:"level"`
	FullPath      bool              `json:"full_path" mapstructure:"full_path"`
	FileWriter    ConfFileWriter    `json:"file_writer" mapstructure:"file_writer"`
	ConsoleWriter ConfConsoleWriter `json:"console_writer" mapstructure:"console_writer"`
	KafKaWriter   ConfKafKaWriter   `json:"kafka_writer" mapstructure:"kafka_writer"`
}

// SetupLog setup log, caller shall use this method
func SetupLog(lc LogConfig) (err error) {
	// global config
	GlobalLevel = getLevel(lc.Level)

	// writer enable
	// 1. if not set level, use global level;
	// 2. if set level, use min level
	validGlobalMinLevel := FATAL // default max level
	validGlobalMinLevelBy := "global"

	fileWriterLevelDefault := GlobalLevel
	consoleWriterLevelDefault := GlobalLevel
	kafkaWriterLevelDefault := GlobalLevel

	if lc.FileWriter.Enable {
		fileWriterLevelDefault = getLevel0(lc.FileWriter.Level, GlobalLevel, WriterNameFile)
		validGlobalMinLevel = util.MinInt(fileWriterLevelDefault, validGlobalMinLevel)
		if validGlobalMinLevel == fileWriterLevelDefault {
			validGlobalMinLevelBy = WriterNameFile
		}
	}

	if lc.ConsoleWriter.Enable {
		consoleWriterLevelDefault = getLevel0(lc.ConsoleWriter.Level, GlobalLevel, WriterNameConsole)
		validGlobalMinLevel = util.MinInt(consoleWriterLevelDefault, validGlobalMinLevel)
		if validGlobalMinLevel == consoleWriterLevelDefault {
			validGlobalMinLevelBy = WriterNameConsole
		}
	}

	if lc.KafKaWriter.Enable {
		kafkaWriterLevelDefault = getLevel0(lc.KafKaWriter.Level, GlobalLevel, WriterNameKafka)
		validGlobalMinLevel = util.MinInt(kafkaWriterLevelDefault, validGlobalMinLevel)
		if validGlobalMinLevel == kafkaWriterLevelDefault {
			validGlobalMinLevelBy = WriterNameKafka
		}
	}

	fullPath := lc.FullPath
	ShowFullPath(fullPath)
	SetLevel(validGlobalMinLevel)

	if lc.FileWriter.Enable {
		w := NewFileWriter()
		w.level = fileWriterLevelDefault
		if err = w.SetPathPattern(lc.FileWriter.LogPath); err != nil {
			return err
		}
		log.Printf("[log4go] enable    " + WriterNameFile + " with level " + LevelFlags[fileWriterLevelDefault])
		Register(w)
	}

	if lc.ConsoleWriter.Enable {
		w := NewConsoleWriter()
		w.level = consoleWriterLevelDefault
		w.SetColor(lc.ConsoleWriter.Color)
		log.Printf("[log4go] enable " + WriterNameConsole + " with level " + LevelFlags[consoleWriterLevelDefault])
		Register(w)
	}

	if lc.KafKaWriter.Enable {
		w := NewKafKaWriter(&lc.KafKaWriter)
		w.level = kafkaWriterLevelDefault
		log.Printf("[log4go] enable   " + WriterNameKafka + " with level " + LevelFlags[kafkaWriterLevelDefault])
		Register(w)
	}

	log.Printf("[log4go] valid global_level(min:%v, flag:%v, by:%v), default(%v, flag:%v)",
		validGlobalMinLevel, LevelFlags[validGlobalMinLevel], validGlobalMinLevelBy, GlobalLevel, LevelFlags[GlobalLevel])
	return nil
}

// SetupLogWithConf setup log with config file
func SetupLogWithConf(file string) (err error) {
	var lc LogConfig
	cnt, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	if err = json.Unmarshal(cnt, &lc); err != nil {
		return
	}
	return SetupLog(lc)
}

func getLevel(flag string) int {
	return getLevel0(flag, DEBUG, "")
}

// default DEBUG
func getLevel0(flag string, defaultFlag int, writer string) int {
	for i, f := range LevelFlags {
		if strings.TrimSpace(strings.ToUpper(flag)) == f {
			return i
		}
	}
	log.Printf("[log4go] no matching level for writer(%v, flag:%v), use default level(%d, flag:%v)", writer, flag, defaultFlag, LevelFlags[defaultFlag])
	return defaultFlag
}
