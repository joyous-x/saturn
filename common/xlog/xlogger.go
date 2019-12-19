package xlog

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

const (
	InvalidLevelMin = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	InvalidLevelMax
)

type XLogLevelDef struct {
	Level   int
	Name    string
	NameFmt string
}

var XLogLevelConf_Name = map[string]*XLogLevelDef{
	"debug": &XLogLevelDef{DebugLevel, "debug", "[Debug]"},
	"info":  &XLogLevelDef{InfoLevel, "info", "[Info ]"},
	"warn":  &XLogLevelDef{WarnLevel, "warn", "[Warn ]"},
	"error": &XLogLevelDef{ErrorLevel, "error", "[Error]"},
	"panic": &XLogLevelDef{ErrorLevel, "panic", "[Panic]"},
}
var XLogLevelConf_Level = map[int]*XLogLevelDef{
	DebugLevel: &XLogLevelDef{DebugLevel, "debug", "[Debug]"},
	InfoLevel:  &XLogLevelDef{InfoLevel, "info", "[Info ]"},
	WarnLevel:  &XLogLevelDef{WarnLevel, "warn", "[Warn ]"},
	ErrorLevel: &XLogLevelDef{ErrorLevel, "error", "[Error]"},
	PanicLevel: &XLogLevelDef{ErrorLevel, "panic", "[Panic]"},
}

///////////////////////////////////////////////////////////////////////////////

type XLogInter interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
}

type XLogConfig struct {
	Path  string `json:path`
	Level int    `json:level`
}

func NewXLogger(config *XLogConfig) *XLogger {
	newer := &XLogger{
		filepath: config.Path,
	}

	newer.writer = newer.rolling()
	newer.level = func(level int) int {
		if level > InvalidLevelMin && level < InvalidLevelMax {
			return level
		}
		return DebugLevel
	}(config.Level)

	if nil == newer.writer {
		return nil
	}

	newer.logger = log.New(newer.writer, "" /*prefix*/, 0 /*log.Lshortfile|log.LstdFlags*/)
	return newer
}

type XLogger struct {
	filepath string
	suffix   string
	level    int
	lock     sync.Mutex
	writer   io.Writer
	logger   *log.Logger
}

func (this *XLogger) rolling() io.Writer {
	if "" == this.filepath {
		return os.Stdout
	}

	this.lock.Lock()
	defer this.lock.Unlock()

	cur_timestamp := time.Now().Format("%y-%m-%d")
	if this.suffix == cur_timestamp {
		return this.writer
	} else {
		f, err := os.OpenFile(this.filepath+"_"+cur_timestamp, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
		if err != nil {
			return this.writer
		}
		this.writer = f
		this.suffix = cur_timestamp
		this.logger.SetOutput(f)
	}

	return this.writer
}

func (this *XLogger) precheck_enable(cur_level int) bool {
	this.rolling()
	return this.level <= cur_level
}

func (this *XLogger) runtime() (string, int) {
	_, file, line, _ := runtime.Caller(4)
	return file, line
}

func (this *XLogger) Print(level int, format string, v ...interface{}) error {
	if false == this.precheck_enable(level) {
		return nil
	}

	level_str := XLogLevelConf_Level[level].NameFmt
	date_str := time.Now().Local().Format("2006-01-02 15:04:05.999999") // ("2006-01-02 15:04:05.999999 -0800 PST") //> time.Now().Local().Format(time.RFC3339Nano)
	file, line := this.runtime()
	msg := fmt.Sprintf(format, v...)

	data := fmt.Sprintf("%v %v %v:%v %v", date_str, level_str, filepath.Base(file), line, msg)
	return this.logger.Output(2, data)
}

func (this *XLogger) Debug(format string, v ...interface{}) {
	this.Print(DebugLevel, format, v...)
}

func (this *XLogger) Info(format string, v ...interface{}) {
	this.Print(InfoLevel, format, v...)
}

func (this *XLogger) Warn(format string, v ...interface{}) {
	this.Print(WarnLevel, format, v...)
}

func (this *XLogger) Error(format string, v ...interface{}) {
	this.Print(ErrorLevel, format, v...)
}

func (this *XLogger) Panic(format string, v ...interface{}) {
	this.Print(PanicLevel, format, v...)
	panic("")
}

///////////////////////////////////////////////////////////////////////////////

var g_config = &XLogConfig{}
var g_XLogger *XLogger
var once sync.Once

func SetCfg(cfg *XLogConfig) error {
	if cfg != nil {
		g_config = cfg
		return nil
	}
	return fmt.Errorf("invalid args")
}

func Inst() *XLogger {
	once.Do(func() {
		g_XLogger = NewXLogger(g_config)
	})
	return g_XLogger
}

func Debug(format string, v ...interface{}) {
	Inst().Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	Inst().Info(format, v...)
}

func Warn(format string, v ...interface{}) {
	Inst().Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	Inst().Error(format, v...)
}

func Panic(format string, v ...interface{}) {
	Inst().Panic(format, v...)
}
