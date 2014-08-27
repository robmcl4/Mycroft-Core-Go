package logging

import (
    "github.com/op/go-logging"
    "os"
    "path/filepath"
)


type Level logging.Level


const (
    CRITICAL Level = Level(logging.CRITICAL)
    ERROR Level    = Level(logging.ERROR)
    WARNING Level  = Level(logging.WARNING)
    NOTICE Level   = Level(logging.NOTICE)
    INFO Level     = Level(logging.INFO)
    DEBUG Level    = Level(logging.DEBUG)
)


var log = logging.MustGetLogger("mycroft")


var (
    Debug    = log.Debug
    Info     = log.Info
    Notice   = log.Notice
    Warning  = log.Warning
    Error    = log.Error
    Critical = log.Critical
    Fatal    = log.Fatal
)


func init() {
    logging.SetLevel(logging.DEBUG, "mycroft")

    format := `%{time:15:04:05.000000} %{level:.4s} %{id:04x} | %{message}`
    logging.SetFormatter(logging.MustStringFormatter(format))

    consoleBackend := logging.NewLogBackend(os.Stderr, "", 0)
    consoleBackend.Color = true
    logging.SetBackend(consoleBackend)
    log.Debug("Console logging enabled")

    // set up the log file
    err := os.MkdirAll("log", os.ModeDir)
    if err != nil {
        log.Error("could not make log directory: %s", err.Error())
        return
    }
    path := filepath.Join("log", "log")
    fileWriter, err := os.OpenFile(path,
                                   os.O_WRONLY|os.O_APPEND|os.O_CREATE,
                                   0660)
    if err != nil {
        log.Error("could not open log file: %s", err.Error())
        return
    }

    fileBackend := logging.NewLogBackend(fileWriter, "", 0)
    fileBackend.Color = false
    logging.SetBackend(logging.MultiLogger(consoleBackend, fileBackend))
    log.Debug("File backend enabled")
}


func SetLevel(l Level) {
    logging.SetLevel(logging.Level(l), "mycroft")
}
