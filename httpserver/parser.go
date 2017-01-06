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

type Stats struct {
	Group string
	Last  Last
}

func (d Day) toString() string {
	return toJson(d)
}

func setStats(group string, last Last) []byte {
	s := Stats{Group: group, Last: last}
	out, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return out
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

func getMax(DataPath string) Last {
	raw, err := ioutil.ReadFile(DataPath + "hour_max.json")
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

func getAllDays(DataPath string) []Day {
	raw, err := ioutil.ReadFile(DataPath + "daily.json")
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

func getDay(ScritpsPath string, args ...string) Day {
	out := RunScript(ScritpsPath+"daily.sh", args...)
	var d Day
	err := json.Unmarshal(out, &d)
	if err != nil {
		panic(err)
	}
	return d
}

func LastHour(ScritpsPath string, DataPath string) Last {

	out := RunScript(ScritpsPath + "last_hour.sh")
	last := getLast(out)
	max := getMax(DataPath)
	fmt.Println("total", last, max)
	if greaterThenMax(last, max) {
		fmt.Println("saving max", last)
		saveMax(last, DataPath)
	}
	return last
}

func saveToday(today Day, DataPath string) {
	days := getAllDays(DataPath)
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
	writeToFile(j, DataPath+"daily.json")
}

func greaterThenMax(last Last, max Last) bool {
	return last.Total > max.Total
}

func saveMax(m Last, DataPath string) {
	j, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		panic(err)
	}
	writeToFile(j, DataPath+"hour_max.json")
}

func writeToFile(data []byte, path string) {

	err := ioutil.WriteFile(path, data, 0777)
	if err != nil {
		panic(err)
	}
}
