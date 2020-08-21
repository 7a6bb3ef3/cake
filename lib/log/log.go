package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"strings"
	"time"
)

// TODO has no rotate io.Writer
const logFileName = "cake.log"

var logLevel zapcore.Level = zap.DebugLevel

var logFile *os.File
var logger *zap.SugaredLogger

func init(){
	f ,e := os.OpenFile(logFileName,os.O_APPEND | os.O_WRONLY | os.O_CREATE ,0755)
	if e != nil{
		panic(e)
	}
	logFile = f
	newLogger()
}

func newLogger() {
	var encfg = zap.NewProductionEncoderConfig()
	//encfg.EncodeTime = zapcore.ISO8601TimeEncoder
	timeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	encfg.EncodeTime = timeEncoder
	//encfg.EncodeLevel = zapcore.CapitalLevelEncoder
	encfg.EncodeLevel = zapcore.CapitalLevelEncoder
	mulw := io.MultiWriter(logFile,os.Stderr)
	ws := zapcore.AddSync(mulw)
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encfg), ws, logLevel)
	logger = zap.New(core).Sugar()
}

func Debug(args ...interface{}){
	logger.Debug(args...)
}

func Info(args ...interface{}){
	logger.Info(args...)
}

func Warn(args ...interface{}){
	logger.Warn(args...)
}

func Error(args ...interface{}){
	logger.Error(args...)
}

func Errorx(msg string ,err error){
	logger.Error(msg ," -> " ,err)
}

func Panic(args ...interface{}){
	logger.Panic(args...)
}

func Fatal(args ...interface{}){
	logger.Fatal(args...)
}

type AdaptLogger struct {
	logger *zap.SugaredLogger
}

func (a *AdaptLogger)Printf(s string ,args ...interface{}){
	if strings.HasSuffix(s ,"\n") {
		s = s[:len(s) - 1]
	}
	a.logger.Infof(s ,args)
}

func GetAdaptLogger() *AdaptLogger{
	return &AdaptLogger{logger: logger}
}
