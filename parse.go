package main

import (
	_ "bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/darkua/babl-dashboard/httpserver"
)

type Day struct {
	Date  string `json:"date"`
	Value int    `json:"value"`
	Error int    `json:"error"`
	L     int    `json:"l"`
	U     int    `json:"u"`
}

func (d Day) toString() string {
	return toJson(d)
}

func toJson(d interface{}) string {
	bytes, err := json.Marshal(d)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return string(bytes)
}

func getAllDays() []Day {
	raw, err := ioutil.ReadFile("./httpserver/static/data/daily.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var d []Day
	json.Unmarshal(raw, &d)
	return d
}

func getDay(args ...string) Day {
	out := RunScript("./scripts/daily.sh", args...)
	var d Day
	json.Unmarshal(out, &d)
	return d
}

func saveToday(today Day) {
	days := getAllDays()
	i := -1
	for k, d := range days {
		if d.Date == today.Date {
			i = k
		}
	}
	if i == -1 {
		days = append(days, today)
	} else {
		days[i] = today
	}
	writeToFile(days, "./httpserver/static/data/daily.json")
}

func writeToFile(days []Day, path string) {
	// buf := make([]byte, 0)
	// out := bytes.NewBuffer(buf)

	j, err := json.MarshalIndent(days, "", "  ")
	if err != nil {
		fmt.Println("err:", err.Error())
	}

	err = ioutil.WriteFile(path, j, 0777)
	if err != nil {
		fmt.Println("err:", err.Error())
	}
}
