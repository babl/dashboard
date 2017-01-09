package httpserver

import (
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

type Counter struct {
	Module            string
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
	r.HandleFunc("/ws/{group}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		group := vars["group"]
		fmt.Println("group", group)
		serveWs(wsHub, w, r, group)
	})

	r.HandleFunc("/lasthour", func(w http.ResponseWriter, r *http.Request) {
		ModuleUser := r.FormValue("user")
		fmt.Println(ModuleUser)
		DataPath := "./httpserver/static/data/" + ModuleUser + "/"
		ScriptsPath := "./scripts/" + ModuleUser + "/"
		last := LastHour(ScriptsPath, DataPath)

		w.Header().Set("Content-Type", "text/plain")
		out := setStats(ModuleUser, last)
		w.Write(out)
	})

	r.HandleFunc("/loyalist", func(w http.ResponseWriter, r *http.Request) {

		t, err := template.ParseFiles(dir + "/loyalist.html") // Parse template file.
		if err != nil {
			panic(err)
		}

		counters := &Counter{
			Module:            "7",
			LastHourDate:      "",
			LastHourReq:       0,
			LastHourErrorRate: 0,
			MaxReq:            0,
			MaxReqDate:        "",
			SuccessPercent:    0,
			ErrorPercent:      0,
		}
		t.Execute(w, *counters) // merge.

	})

	r.HandleFunc("/babl", func(w http.ResponseWriter, r *http.Request) {

		t, err := template.ParseFiles(dir + "/babl.html") // Parse template file.
		if err != nil {
			panic(err)
		}

		counters := &Counter{
			Module:            "all",
			LastHourDate:      "",
			LastHourReq:       0,
			LastHourErrorRate: 0,
			MaxReq:            0,
			MaxReqDate:        "",
			SuccessPercent:    0,
			ErrorPercent:      0,
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
	StartCrons(wsHub, "babl")
	StartCrons(wsHub, "loyalist")
	log.Fatal(srv.ListenAndServe())
}

func StartCrons(wsHub *Hub, ModuleUser string) {
	//setup crons
	c := cron.New()

	DataPath := "./httpserver/static/data/" + ModuleUser + "/"
	ScriptsPath := "./scripts/" + ModuleUser + "/"

	//gather and save daily stats
	c.AddFunc("0 * * * * *", func() {
		t := time.Now()
		today := fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
		d := getDay(ScriptsPath, today)
		saveToday(d, DataPath)
	})

	c.AddFunc("0 * * * * *", func() {
		last := LastHour(ScriptsPath, DataPath)
		out := setStats(ModuleUser, last)
		wsHub.Broadcast <- out //#todo: replace broadcast to all with group channels!
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
