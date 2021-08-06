package main

import (
	"flag"
	"fmt"
	"github.com/comhttp/json/api"
	"github.com/comhttp/json/cfg"
	daemon "github.com/leprosus/golang-daemon"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var log = logrus.New()

func wrapLogger(module string) logrus.FieldLogger {
	return log.WithField("module", module)
}

func parseLogLevel(level string) logrus.Level {
	switch level {
	case "error":
		return logrus.ErrorLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "info", "notice":
		return logrus.InfoLevel
	case "debug":
		return logrus.DebugLevel
	case "trace":
		return logrus.TraceLevel
	default:
		return logrus.InfoLevel
	}
}
func main() {
	// Get cmd line parameters
	command := flag.String("cmd", "", "Command")
	path := flag.String("path", "./", "Path")
	port := flag.String("port", "8080", "Port")
	loglevel := flag.String("loglevel", "info", "Logging level (debug, info, warn, error)")
	flag.Parse()

	log.SetLevel(parseLogLevel(*loglevel))
	cfg.CONFIG = &cfg.Conf{
		Path: *path,
	}
	a, _ := api.NewAPI(*path, nil)

	err := daemon.Init(*command, map[string]interface{}{}, "./daemonized.pid")
	if err != nil {
		return
	}
	switch *command {
	case " start":
		err = daemon.Start()
	case " stop":
		err = daemon.Stop()
	case " restart":
		err = daemon.Stop()
		err = daemon.Start()
	case " status":
		status := "stopped"
		if daemon.IsRun() {
			status = "started"
		}

		fmt.Printf("Application is %s\n", status)

		return
	case "":
	default:
		w := &http.Server{
			Handler:      api.Handler(a),
			Addr:         ":" + *port,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		log.Fatal(w.ListenAndServe())
		fmt.Println("JORM node is on: ", ":"+*port)
	}
}
