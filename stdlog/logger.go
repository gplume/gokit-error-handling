package stdlog

import (
	"log"
	"os"
)

// Log ..
var Log *log.Logger

func init() {
	Log = log.New(os.Stdout, "api: ", log.Lshortfile)
}

// New ...
func New() *log.Logger {
	return Log
}
