package httpserver

import (
	_ "bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	_ "os"
)

type Day struct {
	Date  string `json:"date"`
	Value int    `json:"value"`
	Error int    `json:"error"`
	L     int    `json:"l"`
	U     int    `json:"u"`
}

type Last struct {
	Date  string `json:"date"`
	Total int    `json:"total"`
	Error int    `json:"error"`
}

func (d Day) toString() string {
	return toJson(d)
}

func toJson(d interface{}) string {
	bytes, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func getLast(raw []byte) Last {
	var l Last
	err := json.Unmarshal(raw, &l)
	if err != nil {
		panic(err)
	}
	return l
}

func getMax() Last {
	raw, err := ioutil.ReadFile("./httpserver/static/data/hour_max.json")
	if err != nil {
		panic(err)
	}

	var m Last
	err = json.Unmarshal(raw, &m)
	if err != nil {
		panic(err)
	}
	return m
}

func getAllDays() []Day {
	raw, err := ioutil.ReadFile("./httpserver/static/data/daily.json")
	if err != nil {
		panic(err)
	}

	var d []Day
	err = json.Unmarshal(raw, &d)
	if err != nil {
		panic(err)
	}
	return d
}

func getDay(args ...string) Day {
	out := RunScript("./scripts/daily.sh", args...)
	var d Day
	err := json.Unmarshal(out, &d)
	if err != nil {
		panic(err)
	}
	return d
}

func LastHour() []byte {
	fmt.Println("fodase")
	out := RunScript("./scripts/last_hour.sh")
	last := getLast(out)
	max := getMax()
	fmt.Println("total", last, max)
	if greaterThenMax(last, max) {
		fmt.Println("saving max", last)
		saveMax(last)
	}
	return out
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
	j, err := json.MarshalIndent(days, "", "  ")
	if err != nil {
		panic(err)
	}
	writeToFile(j, "./httpserver/static/data/daily.json")
}

func greaterThenMax(last Last, max Last) bool {
	return last.Total > max.Total
}

func saveMax(m Last) {
	j, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		panic(err)
	}
	writeToFile(j, "./httpserver/static/data/hour_max.json")
}

func writeToFile(data []byte, path string) {

	err := ioutil.WriteFile(path, data, 0777)
	if err != nil {
		panic(err)
	}
}
