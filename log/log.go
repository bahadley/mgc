package log

import (
	//"io/ioutil"
	"log"
	"os"
	//"github.com/bahadley/mgc/config"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func init() {
	/*
		traceOut := ioutil.Discard
		if config.Trace() {
			traceOut = os.Stdout
		}
	*/

	Trace = log.New(os.Stdout,
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
