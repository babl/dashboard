package main

import (
	"fmt"
	_ "net/http"
	"os"
	_ "strings"
	_ "time"

	log "github.com/Sirupsen/logrus"
	. "github.com/darkua/babl-dashboard/httpserver"
)

const Version = "0.0.1"
const clientID = "babl-dashboard"

var debug bool

func main() {
	log.SetOutput(os.Stderr)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.WarnLevel)

	log.Warn("App START")
	run(":8000", false)
}

func run(listen string, dbg bool) {
	debug = dbg
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	// websockets
	log.Warn("App Run Websockets Hub")
	wsHub := NewHub()
	go wsHub.Run()

	// StartHttpServer(listen, wsHub, HandlerRequestHistory, HandlerRequestDetails, HandlerRequestPayload)
	log.Warn(fmt.Sprintf("App Start WebServer running on port %s", listen))
	StartHttpServer(listen, wsHub)

}
