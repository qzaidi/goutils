// The logging package provides common functionality as log rotation, conditional debug logging etc.
package logging

import (
  "os"
  "syscall"
  "flag"
  "log"
  "io/ioutil"
  "os/signal"
)

var stdoutLog string
var stderrLog string
var debugFlag bool

// global logger for debug messages
//  logging.Debug.Println("debug message")
// debug messages are printed only when the program is started with -debug flag
var Debug *log.Logger
var Info *log.Logger

// Init installs the command line options for setting output and error log paths, and exposes
// logging.Debug, which can be used to add code for debug
func init() {
  flag.StringVar(&stdoutLog,"l","","log file for stdout")
  flag.StringVar(&stderrLog,"e","","log file for stderr")
  flag.BoolVar(&debugFlag,"debug",false,"enable debug logging")

  Debug = log.New(ioutil.Discard,"",0)
  Info = log.New(os.Stdout,"info:",log.Ldate|log.Ltime)

  c := make(chan os.Signal, 1)
  signal.Notify(c, syscall.SIGHUP) // listen for sighup
  go sigHandler(c)
}

func sigHandler(c chan os.Signal) {
  // Block until a signal is received.
  for s := range c {
    log.Println("Reloading on :", s)
    LogInit()
  }
}

// App must call LogInit once to setup log redirection
func LogInit() {
  log.Println("Log Init: using ",stdoutLog,stderrLog)
  reopen(1,stdoutLog)
  reopen(2,stderrLog)

  if debugFlag {
    Debug = log.New(os.Stdout,"debug:",log.Lshortfile|log.Ldate|log.Ltime)
    Debug.Println("---- debug mode ----")
  }
}

func reopen(fd int,filename string) {
  if filename == "" {
    return
  }

  logFile,err := os.OpenFile(filename, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)

  if (err != nil) {
    log.Println("Error in opening ",filename,err)
    os.Exit(2)
  }

  syscall.Dup2(int(logFile.Fd()), fd)
}
