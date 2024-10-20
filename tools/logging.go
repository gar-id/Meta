package tools

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	agentCaches "MetaHandler/agent/caches"
	serverCaches "MetaHandler/server/config/caches"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func createLog(loglocation string) *os.File {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// Create folder if doesn't exist
	if _, err := os.Stat(loglocation); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(loglocation, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	loglocation = fmt.Sprintf("%v/meta-handler-agent.log", loglocation)
	logfile, err := os.OpenFile(loglocation, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	return logfile
}

func ZapLogger(storeOption, apps string) *zap.Logger {
	logger := coreZap(storeOption, apps)
	defer logger.Sync()

	return logger
}

func coreZap(storeOption, apps string) *zap.Logger {
	// Setup logging
	var stdout, file zapcore.WriteSyncer
	var logLocation, logLevel string
	switch apps {
	case "agent":
		logLevel = DefaultString(agentCaches.MetaHandlerAgent.AgentConfig.Log.Level, "debug")
		logLocation = DefaultString(agentCaches.MetaHandlerAgent.AgentConfig.Log.Location, path.Join("C", "Program Files", "Meta Handler", "logs"))
	case "server":
		logLevel = DefaultString(serverCaches.MetaHandlerServer.MetaHandlerServer.Log.Level, "debug")
		logLocation = DefaultString(serverCaches.MetaHandlerServer.MetaHandlerServer.Log.Location, path.Join("/", "var", "log", "Meta Handler"))

	}

	if storeOption == "file" {
		logfile := createLog(logLocation)
		file = zapcore.AddSync(logfile)
	} else if storeOption == "console" {
		stdout = zapcore.AddSync(os.Stdout)
	} else if storeOption == "both" {
		logfile := createLog(logLocation)
		file = zapcore.AddSync(logfile)
		stdout = zapcore.AddSync(os.Stdout)
	} else {
		log.Fatal("storeOption is undefined or using unknown string.")
	}

	// Log level
	var level zap.AtomicLevel
	switch logLevel {
	case "debug":
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warning":
		level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "panic":
		level = zap.NewAtomicLevelAt(zap.PanicLevel)
	case "fatal":
		level = zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.FunctionKey = "func"
	productionCfg.EncodeDuration = zapcore.MillisDurationEncoder
	productionCfg.EncodeName = zapcore.FullNameEncoder
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	productionCfg.EncodeCaller = zapcore.FullCallerEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	developmentCfg.CallerKey = ""
	developmentCfg.EncodeCaller = nil

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	var core zapcore.Core
	if storeOption == "file" {
		core = zapcore.NewTee(
			zapcore.NewCore(fileEncoder, file, level),
		)
	} else if storeOption == "console" {
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, stdout, level),
		)
	} else if storeOption == "both" {
		core = zapcore.NewTee(
			zapcore.NewCore(fileEncoder, file, level),
			zapcore.NewCore(consoleEncoder, stdout, level),
		)
	}

	return zap.New(core, zap.AddCaller(), zap.WithCaller(true))
}
