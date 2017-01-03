package httpserver

import (
	"html/template"
	"net/http"
	"os"
	_ "strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	. "github.com/larskluge/babl-server/utils"
)

// func GetVarsBlockSize(r *http.Request, defaultvalue int64) int64 {
// 	result := defaultvalue
// 	vars := mux.Vars(r)
// 	blocksize := vars["blocksize"]
// 	if blocksize != "" {
// 		bsize, errParse := strconv.ParseInt(blocksize, 10, 64)
// 		if errParse == nil {
// 			result = bsize
// 		}
// 	}
// 	return result
// }

type Counter struct {
	Module int
}

func StartHttpServer(listen string, wsHub *Hub) {

	pwd, err := os.Getwd()
	Check(err)
	dir := pwd + "/httpserver/static"
	//fmt.Println("WorkingDir: ", pwd)
	//fmt.Println("HttpServer: ", dir)
	r := mux.NewRouter()

	// REST API
	// r.HandleFunc("/api/request/history", HandlerRequestHistory).Methods("GET").Queries("blocksize", "{blocksize}")
	// r.HandleFunc("/api/request/history", HandlerRequestHistory).Methods("GET")
	// r.HandleFunc("/api/request/details/{requestid:.*}", HandlerRequestDetails).Methods("GET")
	// r.HandleFunc("/api/request/payload/{topic:.*}/{partition:[0-9]+}/{offset:[0-9]+}", HandlerRequestPayload).Methods("GET")

	// websockets
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(wsHub, w, r)
	})

	r.HandleFunc("/demo", func(w http.ResponseWriter, r *http.Request) {
		wsHub.Broadcast <- []byte("Hello World!!!")
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("demo!"))
	})

	r.HandleFunc("/loyalist", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles(dir + "/loyalist.html") // Parse template file.
		if err != nil {
			panic(err)
		}
		counters := &Counter{Module: 7}
		t.Execute(w, *counters) // merge.

	})

	// Static files and assets
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(dir))))

	srv := &http.Server{
		Handler:      r,
		Addr:         listen,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
