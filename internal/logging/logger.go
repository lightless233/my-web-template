package logging

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// var logger zap.Logger
var sugarLogger zap.SugaredLogger

// 暂时不用这个 logger，使用下面的 sugar logger
//func Logger() *zap.Logger {
//	return &logger
//}

func Sugar() *zap.SugaredLogger {
	return &sugarLogger
}

const LogDirName = "logs"
const DefaultLogFilename = "app.log"

// GetExecPath 获取可执行文件的路径，默认将日志文件放到可执行文件同级，如果有其他需求再修改这部分代码
func GetExecPath() string {
	file, _ := exec.LookPath(os.Args[0])
	execPath, _ := filepath.Abs(file)
	execPath = execPath[:strings.LastIndex(execPath, string(os.PathSeparator))]
	return execPath
}

func InitLogger(debug bool, logFilename string) error {
	execPath := GetExecPath()
	logDir := path.Join(execPath, LogDirName)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.Mkdir(logDir, os.ModePerm); err != nil {
			return err
		}
	}

	if logFilename == "" {
		logFilename = DefaultLogFilename
	}

	lumberjackLogger := &lumberjack.Logger{
		Filename:   path.Join(logDir, logFilename),
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     28,
		Compress:   false,
	}

	// 基础编码配置
	baseEncoderConfig := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "name",
		CallerKey:        "caller",
		FunctionKey:      "function",
		MessageKey:       "message",
		StacktraceKey:    zapcore.OmitKey,
		ConsoleSeparator: "|",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.CapitalLevelEncoder, // 默认不带颜色
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(time.RFC3339))
		},
		EncodeName:     zapcore.FullNameEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 控制台编码配置（带颜色）
	consoleEncoderConfig := baseEncoderConfig
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 控制台始终带颜色
	if debug {
		consoleEncoderConfig.ConsoleSeparator = " " // 调试模式下使用空格分隔符
	}

	// 文件编码配置（不带颜色）
	fileEncoderConfig := baseEncoderConfig
	fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 文件不带颜色

	// 创建编码器
	fileEncoder := zapcore.NewConsoleEncoder(fileEncoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

	// 创建输出目标
	fileSyncer := zapcore.AddSync(lumberjackLogger)
	consoleSyncer := zapcore.AddSync(os.Stdout)

	// 创建核心（Core）
	fileCore := zapcore.NewCore(fileEncoder, fileSyncer, zapcore.DebugLevel)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleSyncer, zapcore.DebugLevel)

	// 合并核心
	core := zapcore.NewTee(fileCore, consoleCore)

	// 创建Logger
	logger := zap.New(core, zap.AddCaller())
	sugarLogger = *logger.Sugar()

	return nil
}
