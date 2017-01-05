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

func WSBroadcast(wsHub *Hub, chData chan *[]byte) {
	for msg := range chData {
		fmt.Printf("%s\n", *msg)
		wsHub.Broadcast <- *msg
	}
}

func run(listen string, dbg bool) {
	debug = dbg
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	// websockets
	log.Warn("App Run Websockets Hub")
	wsHub := NewHub()
	chData := make(chan *[]byte)
	go wsHub.Run()
	go WSBroadcast(wsHub, chData)

	// other higher level go rotines go here
	// log.Warn("App Save/Broadcast Data")
	// go MonitorRequest(chQAData, chHistory, chWSHistory, chDetails)
	// go SaveRequestHistory(s.kafkaProducer, kafkaTopicHistory, chHistory)

	// go SaveRequestDetails(s.kafkaProducer, kafkaTopicDetails, chDetails, cacheDetails)

	// // http callback function handler for Request History
	// // $ http 127.0.0.1:8888/api/request/history
	// // $ http 127.0.0.1:8888/api/request/history?blocksize=20
	// HandlerRequestHistory := func(w http.ResponseWriter, r *http.Request) {
	// 	lastn := GetVarsBlockSize(r, 10)
	// 	rhJson := ReadRequestHistory(s.kafkaClient, kafkaTopicHistory, lastn)
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.Write(rhJson)
	// }
	// // http callback function handler for Request Details
	// // http 127.0.0.1:8888/api/request/details/12345
	// HandlerRequestDetails := func(w http.ResponseWriter, r *http.Request) {
	// 	vars := mux.Vars(r)
	// 	rhJson := ReadRequestDetailsFromCache(vars["requestid"], cacheDetails)
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.Write(rhJson)
	// }
	// // http callback function handler for Request Payload
	// // http 127.0.0.1:8888/api/request/payload/babl.babl.Events.IO/4/3460
	// HandlerRequestPayload := func(w http.ResponseWriter, r *http.Request) {
	// 	vars := mux.Vars(r)
	// 	rhJson := ReadRequestPayload(s.kafkaClient, vars["topic"], vars["partition"], vars["offset"])
	// 	w.Header().Set("Content-Type", "application/octet-stream")
	// 	w.Write(rhJson)
	// }

	// StartHttpServer(listen, wsHub, HandlerRequestHistory, HandlerRequestDetails, HandlerRequestPayload)
	log.Warn(fmt.Sprintf("App Start WebServer running on port %s", listen))
	StartHttpServer(listen, wsHub)

}
