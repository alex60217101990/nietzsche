package logger

import (
	"log"
	"os"
	"sync"

	"github.com/fatih/color"
)

var (
	PackageOnceLoad sync.Once

	CmdServer = color.New(color.FgHiGreen, color.Bold)
	CmdError  = color.New(color.FgHiRed, color.Bold)
	CmdInfo   = color.New(color.FgHiBlue, color.Bold)

	AppLogger  Logger
	RepoLogger *log.Logger
)

func InitLoggerSettings() {
	PackageOnceLoad.Do(func() {
		CmdServer.Println("Run once - 'logger' package loading...")

		RepoLogger = log.New(os.Stderr, "[REPO] ", log.LstdFlags|log.Lshortfile|log.LUTC)
		AppLogger = NewZapLogger()
	})
}

func CloseLoggers() {
	AppLogger.Close()

	CmdServer.Println("'logger' package stoped...")
}
