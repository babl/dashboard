// package httpserver

// import (

// )

// func (s *MainSuite) TestNoEndpoints(c *C) {
// 	c.Assert(ParseEndpoints(""), DeepEquals, []rr.Endpoint{})
// }

// func (s *MainSuite) TestOneEndpoint(c *C) {
// 	c.Assert(ParseEndpoints("babl.sh"), DeepEquals, []rr.Endpoint{rr.Endpoint{"babl.sh", 1}})
// }

// func (s *MainSuite) TestTwoEndpoints(c *C) {
// 	c.Assert(ParseEndpoints("babl.sh,v5.babl.sh"), DeepEquals, []rr.Endpoint{rr.Endpoint{"babl.sh", 1}, rr.Endpoint{"v5.babl.sh", 1}})
// }

// func (s *MainSuite) TestWeightedEndpoints(c *C) {
// 	c.Assert(ParseEndpoints("babl.sh;q=9, v5.babl.sh;q=1"), DeepEquals, []rr.Endpoint{rr.Endpoint{"babl.sh", 9}, rr.Endpoint{"v5.babl.sh", 1}})
// }

package httpserver

import (
	"encoding/json"
	"flag"
	"os"
	"testing"
)

func saveOneDay(day string, path string) {

	out := []byte("{\"date\": \"" + day + "\", \"value\": 1, \"error\": 1, \"l\": 1, \"u\": 1}")
	var item Day
	json.Unmarshal(out, &item)
	saveToday(item, path)
}

func expectDayToBeSaved(day string, path string) int {
	//expect the item to be saved in file
	data := getAllDays(path)
	i := -1
	for k, d := range data {
		if d.Date == day {
			i = k
		}
	}
	return i
}
func expectDayToBeUnique(day string, path string) int {
	//expect the item to be saved in file
	data := getAllDays(path)
	i := 0
	for _, d := range data {
		if d.Date == day {
			i++
		}
	}
	return i
}

func TestSaveToday(t *testing.T) {
	day := "2011-01-01"
	path := "./static/test/"
	saveOneDay(day, path)
	i := expectDayToBeSaved(day, path)
	if i == -1 {
		t.Fatalf("data item not saved correctly %s", day)
	}
}

func TestSaveAnyDays(t *testing.T) {
	day1 := "2011-01-02"
	path := "./static/test/"
	saveOneDay(day1, path)
	i := expectDayToBeSaved(day1, path)
	if i == -1 {
		t.Fatalf("data item not saved correctly %s", day1)
	}

	day2 := "2011-01-03"
	saveOneDay(day2, path)
	i = expectDayToBeSaved(day2, path)
	if i == -1 {
		t.Fatalf("data item not saved correctly %s", day2)
	}
}
func TestUpdateDay(t *testing.T) {
	day := "2011-01-04"
	path := "./static/test/"
	saveOneDay(day, path)
	saveOneDay(day, path)
	i := expectDayToBeUnique(day, path)
	if i != 1 {
		t.Fatalf("data item not unique  %s", day)
	}
}

func TestMain(m *testing.M) {
	flag.Parse()

	path := "./static/test/daily.json"
	writeToFile([]byte("[]"), path)

	os.Exit(m.Run())
}
