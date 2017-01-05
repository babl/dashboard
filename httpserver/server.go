package httpserver

import (
	// "fmt"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	_ "strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	. "github.com/larskluge/babl-server/utils"
	"github.com/robfig/cron"
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
	Module            int
	LastHourDate      string
	LastHourReq       int
	LastHourError     int
	LastHourErrorRate float32
	MaxReq            int
	MaxReqDate        string
	SuccessPercent    float32
	ErrorPercent      float32
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

	r.HandleFunc("/lasthour", func(w http.ResponseWriter, r *http.Request) {

		out := LastHour()
		fmt.Println("refresh", string(out))
		w.Header().Set("Content-Type", "text/plain")
		w.Write(out)
	})

	r.HandleFunc("/loyalist", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles(dir + "/loyalist.html") // Parse template file.
		if err != nil {
			panic(err)
		}

		counters := &Counter{
			Module:            7,
			LastHourDate:      "2017-01-03 16:35",
			LastHourReq:       761,
			LastHourErrorRate: 5.64,
			MaxReq:            1000,
			MaxReqDate:        "2017-01-03 15:35",
			SuccessPercent:    76.1,
			ErrorPercent:      5.64,
		}
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

	//setup crons
	StartCrons(wsHub)
	log.Fatal(srv.ListenAndServe())
}

func StartCrons(wsHub *Hub) {
	//setup crons
	c := cron.New()

	//gather and save daily stats
	c.AddFunc("0 * * * * *", func() {
		t := time.Now()
		today := fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
		d := getDay(today)
		saveToday(d)
	})

	c.AddFunc("0 * * * * *", func() {
		out := LastHour()
		wsHub.Broadcast <- out
	})

	c.Start()
}

func RunScript(cmd string, args ...string) []byte {
	output, err := exec.Command(cmd, args...).Output()
	if err != nil {
		log.Fatal(err)
	}
	return output
}
