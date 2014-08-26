package logging

import (
    "github.com/op/go-logging"
    "os"
    "path/filepath"
)


type Level logging.Level


const (
    CRITICAL Level = Level(logging.CRITICAL)
    ERROR Level = Level(logging.ERROR)
    WARNING Level = Level(logging.WARNING)
    NOTICE Level = Level(logging.NOTICE)
    INFO Level = Level(logging.INFO)
    DEBUG Level = Level(logging.DEBUG)
)


var log = logging.MustGetLogger("mycroft")
var Debug = log.Debug
var Info  = log.Info
var Notice = log.Notice
var Warning = log.Warning
var Error = log.Error
var Critical = log.Critical
var Fatal = log.Fatal


func SetupLogging(level Level) {
    logging.SetLevel(logging.Level(level), "mycroft")

    format := `%{color}%{time:15:04:05.000000} %{level:.4s} %{id:04x}%{color:reset} | %{message}`
    logging.SetFormatter(logging.MustStringFormatter(format))

    consoleBackend := logging.NewLogBackend(os.Stderr, "", 0)
    logging.SetBackend(consoleBackend)
    log.Debug("Console logging enabled")

    // set up the log file
    err := os.MkdirAll("log", os.ModeDir)
    if err != nil {
        log.Error("could not make log directory: %s", err.Error())
        return
    }
    path := filepath.Join("log", "log")
    file, err := os.OpenFile(path,
                             os.O_WRONLY|os.O_APPEND|os.O_CREATE,
                             0660)
    if err != nil {
        log.Error("could not open log file: %s", err.Error())
        return
    }

    format = `%{time:15:04:05.000000} %{level:.4s} %{id:04x} | %{message}`
    logging.SetFormatter(logging.MustStringFormatter(format))
    fileBackend := logging.NewLogBackend(file, "", 0)

    logging.MultiLogger(fileBackend, consoleBackend)
    log.Debug("File backend enabled")
}


func SetLevel(l Level) {
    logging.SetLevel(logging.Level(l), "mycroft")
}
