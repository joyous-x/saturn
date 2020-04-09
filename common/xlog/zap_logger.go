package xlog

import (
	"encoding/json"
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

var g_logger *ZapWriter

func InitZapLogger(logInConsole bool) {
	g_logger = initLogger(logInConsole)
}

func initLogger(logInConsole bool) *ZapWriter {
	var zapWriter ZapWriter
	logName := fmt.Sprintf("%s_%v", "/data/apps/zab", rand.Int63n(10000))
	zapWriter.Init(logName, zap.InfoLevel, 100, 30, 28, true, logInConsole)
	Debug("initZapLogger ok: path=%#v", logName)
	return &zapWriter
}

func RandHit(rate int64) bool {
	return (rand.Int63n(100) < rate)
}

func AsyncZapWriteJson(strJson string) {
	g_logger.AsyncWriteJson(strJson)
}

func AsyncZapWriteMap(data map[string]interface{}) {
	g_logger.AsyncWriteMap(data)
}

////////////////////////////////////////////////////////

type ZapWriter struct {
	logger      *zap.Logger
	bufferJson  chan string
	bufferMap   chan map[string]interface{}
	quit        chan os.Signal
	probability int
}

func (z *ZapWriter) Init(logPath string, level zapcore.Level, maxSize, maxBackups, maxAge int, jsonFormat, logInConsole bool) error {
	hook := lumberjack.Logger{
		Filename:   logPath,    // 日志文件路径
		MaxSize:    maxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: maxBackups, // 日志文件最多保存多少个备份
		MaxAge:     maxAge,     // 文件最多保存多少天
		Compress:   false,      // 是否压缩
	}
	var syncer zapcore.WriteSyncer
	if logInConsole {
		syncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook))
	} else {
		syncer = zapcore.AddSync(&hook)
	}

	var encoder zapcore.Encoder
	if jsonFormat {
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	} else {
		encoder = zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
	}

	core := zapcore.NewCore(
		encoder, // 编码器配置
		syncer,  // 打印到控制台和文件
		level,   // 日志级别
	)

	z.bufferMap = make(chan map[string]interface{}, 1000)
	z.bufferJson = make(chan string, 600)
	z.logger = zap.New(core)
	z.quit = make(chan os.Signal, 1)
	signal.Notify(z.quit, os.Interrupt)
	go z.run()
	return nil
}

func (z *ZapWriter) WriteMap(flag string, data map[string]interface{}) error {
	fields := make([]zap.Field, 0)
	fields = append(fields, zap.Any("date", time.Now().Format("2006-01-02 15:04:05")))
	for k, v := range data {
		f := zap.Any(k, v)
		fields = append(fields, f)
	}
	z.logger.Info(flag, fields...)
	return nil
}

func (z *ZapWriter) WriteJson(flag, strJson string) error {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(strJson), &data)
	if err != nil {
		Warn("ZapWriter.Write isn't json: %v", strJson)
		data = map[string]interface{}{
			"unknown": strJson,
		}
	}
	return z.WriteMap(flag, data)
}

func (z *ZapWriter) AsyncWriteJson(data string) error {
	select {
	case z.bufferJson <- data:
		Debug("AsyncWriteJson ok: str=%#v", data)
	case <-time.After(time.Millisecond * 300):
		Debug("AsyncWriteJson overtime: str=%#v", data)
	}
	return nil
}

func (z *ZapWriter) AsyncWriteMap(data map[string]interface{}) error {
	select {
	case z.bufferMap <- data:
		Debug("AsyncWriteMap ok: map=%#v", data)
	case <-time.After(time.Millisecond * 300):
		Debug("AsyncWriteMap overtime: map=%#v", data)
	}
	return nil
}

func (z *ZapWriter) run() {
	for {
		select {
		case v := <-z.bufferJson:
			z.WriteJson("json", v)
		case m := <-z.bufferMap:
			z.WriteMap("map", m)
		case <-z.quit:
			Warn("zap writer is interrupted")
			return
		}
	}
}
