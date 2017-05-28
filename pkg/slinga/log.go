package slinga

import (
	log "github.com/Sirupsen/logrus"
	"fmt"
	"os"
)

var tracing *ScreenLogger
var debug *log.Logger

// ScreenLogger contains is a logger that prints onto the screen and supports on/off
type ScreenLogger struct {
	enabled  bool
}

func (logger *ScreenLogger) setEnable(enabled bool) {
	logger.enabled = enabled
}

func (logger *ScreenLogger) Printf(depth int, format string, args ...interface{}) {
	if logger.enabled {
		indent := ""
		for n := 0; n <= 4*depth; n++ {
			indent = indent + " "
		}
		format = indent + format + "\n"
		fmt.Printf(format, args...)
	}
}

func (logger *ScreenLogger) Println() {
	if logger.enabled {
		fmt.Println()
	}
}

func SetDebugLevel(level log.Level) {
	debug.Level = level
}

func init() {
	tracing = &ScreenLogger{}

	debug = log.New()
	debug.Out, _ = os.OpenFile(GetAptomiDBDir() + "/" + "debug.log", os.O_CREATE|os.O_WRONLY, 0644)

	// Don't log much by default. It will be overridden with "--debug" from CLI
	debug.Level = log.PanicLevel

	// Add a hook to print important errors to stdout as well
	debug.Hooks.Add(&LogHook{})
}

type LogHook struct {

}

func (l *LogHook) Levels() []log.Level {
	return []log.Level{
		log.WarnLevel,
		log.ErrorLevel,
		log.FatalLevel,
		log.PanicLevel,
	}
}

func (l *LogHook) Fire(e *log.Entry) error {
	fmt.Println("Error!")
	fmt.Printf("  %s\n", e.Message)
	fmt.Printf("  %v\n", e.Data)
	return nil
}