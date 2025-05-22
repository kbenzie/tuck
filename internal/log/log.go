package log

import (
	"io"
	"log"
	"os"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var (
	logger *log.Logger
	Level  LogLevel
)

func Debugln(args ...any) {
	if Level <= LevelDebug {
		logger.Println(append([]any{"debug:"}, args...)...)
	}
}

func Debugf(format string, args ...any) {
	if Level <= LevelDebug {
		logger.Printf("debug: "+format, args...)
	}
}

func Infoln(args ...any) {
	if Level <= LevelInfo {
		logger.Println(append([]any{"info:"}, args...)...)
	}
}

func Infof(format string, args ...any) {
	if Level <= LevelInfo {
		logger.Printf("info: "+format, args...)
	}
}

func Warnln(args ...any) {
	if Level <= LevelWarn {
		logger.Println(append([]any{"warn:"}, args...)...)
	}
}

func Warnf(format string, args ...any) {
	if Level <= LevelWarn {
		logger.Printf("warn: "+format, args...)
	}
}

func Errorln(args ...any) {
	if Level <= LevelError {
		logger.Println(append([]any{"error:"}, args...)...)
	}
}

func Errorf(format string, args ...any) {
	if Level <= LevelError {
		logger.Printf("error: "+format, args...)
	}
}

func Fatalln(args ...any) {
	if Level <= LevelFatal {
		logger.Fatalln(append([]any{"fatal:"}, args...)...)
	}
}

func Fatalf(format string, args ...any) {
	if Level <= LevelFatal {
		logger.Fatalf("fatal: "+format, args...)
	}
}

func SetOutput(out io.Writer, flag int) {
	logger = log.New(out, "tuck ", flag)
}

func SetLevel(level LogLevel) {
	Level = level
}

func init() {
	SetOutput(os.Stderr, 0)
	SetLevel(LevelInfo)
}
