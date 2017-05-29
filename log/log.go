package log

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/bahadley/mgc/system"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func init() {
	traceOut := ioutil.Discard
	if system.Trace() {
		traceOut = os.Stdout
	}

	Trace = log.New(traceOut,
		"TRACE: ",
		log.Ldate|log.Lmicroseconds|log.Lshortfile)

	Info = log.New(os.Stderr,
		"INFO: ",
		log.Ldate|log.Ltime)

	Warning = log.New(os.Stderr,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(os.Stderr,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
