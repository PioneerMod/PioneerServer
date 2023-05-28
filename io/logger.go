package io

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/rivo/tview"
)

type LogLevel int

const (
	Debug   LogLevel = 0
	Info    LogLevel = 1
	Success LogLevel = 2
	Warn    LogLevel = 3
	Error   LogLevel = 4
	Fatal   LogLevel = 5
)

type logWriter struct {
	view *tview.TextView
	mu   sync.Mutex
}

func (w *logWriter) Write(str string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	fmt.Fprint(w.view, str)
}

var logLevel LogLevel = Info
var timeFormat string = "2006-01-02 15:04:05"
var writer *logWriter

func InitLogger(config LoggerConfig, logView *tview.TextView) {
	logLevel = LogLevel(config.Level)
	timeFormat = config.TimeFormat

	writer = &logWriter{
		view: logView,
		mu:   sync.Mutex{},
	}
}

func Log(level LogLevel, rawMessage string) {
	if level < logLevel {
		return
	}

	time := getTime()
	message := fmt.Sprintf("[grey]%s[-:-:-] %s [cyan]>[-:-:-] [white]%s[-:-:-]\n", time, getLogLevelPrefix(level), rawMessage)

	writer.Write(message)
}

func Logf(level LogLevel, format string, a ...interface{}) {
	if level < logLevel {
		return
	}

	time := getTime()
	message := fmt.Sprintf("[grey]%s[-:-:-] %s [cyan]>[-:-:-] [white]%s[-:-:-]\n", time, getLogLevelPrefix(level), fmt.Sprintf(format, a...))

	writer.Write(message)
}

func getTime() string {
	return time.Now().Format(timeFormat)
}

func getLogLevelPrefix(level LogLevel) string {
	switch level {
	case Debug:
		return "[pink]DBG[-:-:-]"
	case Info:
		return "[blue]INF[-:-:-]"
	case Success:
		return "[green]SUC[-:-:-]"
	case Warn:
		return "[yellow]WRN[-:-:-]"
	case Error:
		return "[red]ERR[-:-:-]"
	case Fatal:
		return "[red]FTL[-:-:-]"
	default:
		return ""
	}
}

func PrintSplashScreen(w *tview.TextView) {
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, `        ____  _                            __  ___          __
	   / __ \(_)___  ____  ___  ___  _____/  |/  /___  ____/ /
	  / /_/ / / __ \/ __ \/ _ \/ _ \/ ___/ /|_/ / __ \/ __  / 
	 / ____/ / /_/ / / / /  __/  __/ /  / /  / / /_/ / /_/ /  
	/_/   /_/\____/_/ /_/\___/\___/_/  /_/  /_/\____/\__,_/   
	`)
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "\tPioneer Server v"+os.Getenv("PIONEER_SRV_VERS")+" by DasDarki & Pythagorion\n")
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "\n")
}
