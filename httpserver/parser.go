package httpserver

import (
	"encoding/json"
	"io/ioutil"
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

type Modules struct {
	Date string `json:"date"`
	Data []struct {
		Module string `json:"module"`
		Data   struct {
			Value int `json:"value"`
			Error int `json:"error"`
			L     int `json:"l"`
			U     int `json:"u"`
		} `json:"data"`
	} `json:"data"`
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

func getAllDaysModules(DataPath string) []Modules {
	raw, err := ioutil.ReadFile(DataPath + "modules_daily.json")
	if err != nil {
		panic(err)
	}

	var d []Modules
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

func getModuleData(ScritpsPath string, args ...string) Modules {
	out := RunScript(ScritpsPath+"modules.sh", args...)
	var m Modules
	err := json.Unmarshal(out, &m)
	if err != nil {
		panic(err)
	}
	return m
}

func LastHour(ScritpsPath string, DataPath string) Last {

	out := RunScript(ScritpsPath + "last_hour.sh")
	last := getLast(out)
	max := getMax(DataPath)
	if greaterThenMax(last, max) {
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
	j = nil
}

func saveTodayModules(today Modules, DataPath string) {
	days := getAllDaysModules(DataPath)
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
	writeToFile(j, DataPath+"modules_daily.json")
	j = nil
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
	j = nil
}

func writeToFile(data []byte, path string) {

	err := ioutil.WriteFile(path, data, 0777)
	if err != nil {
		panic(err)
	}
}
