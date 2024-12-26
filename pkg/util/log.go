package util

import (
    "sync"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
)

var (
    loggerInstance *zap.Logger
    once           sync.Once
)

func GetLogger() *zap.Logger {
    once.Do(func() {
        lumberjackLogger := &lumberjack.Logger{
            Filename:   "debug.log", 
            MaxSize:    256,                       
            MaxBackups: 1,                         
            MaxAge:     28,                        
            Compress:   true,                      
        }

        encoderConfig := zap.NewProductionEncoderConfig()
        encoderConfig.TimeKey = "time"
        encoderConfig.LevelKey = "level"
        encoderConfig.NameKey = "logger"
        encoderConfig.CallerKey = "caller"
        encoderConfig.MessageKey = "msg"
        encoderConfig.StacktraceKey = "stacktrace"
        encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
        encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

        core := zapcore.NewCore(
            zapcore.NewJSONEncoder(encoderConfig),                  
            zapcore.AddSync(lumberjackLogger),                      
            zap.InfoLevel,                                          
        )

        loggerInstance = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
    })

    return loggerInstance
}

