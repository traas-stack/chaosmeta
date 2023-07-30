package log

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io/ioutil"
	"log"
	"os"
	"sigs.k8s.io/yaml"
	"strings"
	"sync"
)

const TimeLayout = "2006-01-02 15:04:05"

var DefaultLogger *LoggerStruct

type LogOption struct {
	LogPath    string // Log路径
	MaxAge     int    // 日志保留天数
	MaxSize    int    // 日志大小上限(MB)
	MaxBackups int    // 日志备份数量
	Level      string // 日志等级(fatal,panic,error,warn,info,debug)
	OutPutType string // 日志输出模式(FileOutPut,StdErrPut,BothFileAndStdErrPut,NoOutPut)
}

type LoggerStruct struct {
	level        zapcore.Level
	depth        int
	globalMutex  sync.RWMutex
	globalLogger *zap.Logger
	pool         buffer.Pool
	logOption    LogOption
}

func init() {
	var err error
	DefaultLogger, err = IntiLoggerSizeRotate(LogOption{
		LogPath:    "./default.log",
		Level:      "info",
		OutPutType: "bothfileandstderrput",
	})
	if err != nil {
		panic("init default logger failed")
	}
}

// NewWithConfigPath
// @Description: initialize logger from configuration file
// @param path: config path
// @return *LoggerStruct: logger object
// @return error
func NewWithConfigPath(path string) (*LoggerStruct, error) {
	if len(path) == 0 {
		return nil, errors.New("config is empty.")
	}
	// 解析配置文件
	var logOption LogOption
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("read config failed")
	}

	if err := yaml.Unmarshal(data, &logOption); err != nil {
		log.Fatalf("failed to parse config file, error: %s", err)
	}

	if !strings.HasSuffix(logOption.LogPath, ".log") {
		logOption.LogPath = logOption.LogPath + ".log"
	}

	return IntiLoggerSizeRotate(logOption)
}

// getZapLogLevel
// @Description: read the log level from the configuration and convert it to the zap log level
// @param option: configuration structure
// @return zapcore.Level: the zap log level
func getZapLogLevel(option LogOption) zapcore.Level {
	switch strings.ToLower(option.Level) {
	case "fatal":
		return zap.FatalLevel
	case "panic":
		return zap.PanicLevel
	case "error":
		return zap.ErrorLevel
	case "warn":
		return zap.WarnLevel
	case "info":
		return zap.InfoLevel
	case "debug":
		return zap.DebugLevel
	default:
		return zap.DebugLevel
	}
}

// getOutPutPaths
// @Description: read the log output path from the configuration file
// @param option: configuration structure
// @return []string: slice of out put paths
func getOutPutPaths(option LogOption, hook lumberjack.Logger) []zapcore.WriteSyncer {
	switch strings.ToLower(option.OutPutType) {
	case "fileoutput":
		return []zapcore.WriteSyncer{zapcore.AddSync(&hook)}
	case "stderrput":
		return []zapcore.WriteSyncer{zapcore.AddSync(os.Stderr)}
	case "bothfileandstderrput":
		return []zapcore.WriteSyncer{zapcore.AddSync(os.Stderr), zapcore.AddSync(&hook)}
	case "nooutput":
		return nil
	default:
		return []zapcore.WriteSyncer{zapcore.AddSync(os.Stderr), zapcore.AddSync(&hook)}
	}
}

// IntiLoggerSizeRotate
// @Description: initialize log according to file size
// @param option: configuration structure
// @return *LoggerStruct: logger object
// @return error
func IntiLoggerSizeRotate(option LogOption) (*LoggerStruct, error) {
	var l LoggerStruct
	l.pool = buffer.NewPool()

	hook := lumberjack.Logger{
		Filename:   option.LogPath,
		MaxSize:    option.MaxSize,
		MaxBackups: option.MaxBackups,
		MaxAge:     option.MaxAge,
		Compress:   true,
	}

	if hook.Filename == "" {
		hook.Filename = "/home/admin/logs/default.log"
	}
	if hook.MaxBackups == 0 {
		hook.MaxBackups = 3
	}
	if hook.MaxSize == 0 {
		hook.MaxSize = 128
	}
	if hook.MaxAge == 0 {
		hook.MaxAge = 7
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime:    zapcore.TimeEncoderOfLayout(TimeLayout),
		//EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),                     // 编码器配置
		zapcore.NewMultiWriteSyncer(getOutPutPaths(option, hook)...), // 打印到控制台和文件
		getZapLogLevel(option),                                       // 日志级别
	)

	caller := zap.AddCaller()
	development := zap.Development()
	logger := zap.New(core, caller, development)

	if logger == nil {
		panic("log 初始化失败")
	}
	l.level = getZapLogLevel(option)
	l.globalLogger = logger
	l.logOption = option

	return &l, nil
}

func (l *LoggerStruct) getLogOption() LogOption {
	return l.logOption
}

func (l *LoggerStruct) setLogOption(logOption LogOption) {
	l.logOption = logOption
}
