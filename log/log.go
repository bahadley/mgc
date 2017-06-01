package log

import (
	"io/ioutil"
	"log"
	"os"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func SetTrace(trace bool) {
	traceOut := ioutil.Discard
	if trace {
		traceOut = os.Stdout
	}

	Trace = log.New(traceOut,
		"TRACE: ",
		log.Ldate|log.Lmicroseconds|log.Lshortfile)
}

func init() {
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
