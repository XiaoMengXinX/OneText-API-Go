package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type onetextData []struct {
	Text string   `json:"text"`
	By   string   `json:"by"`
	From string   `json:"from"`
	Time []string `json:"time"`
	URI  string   `json:"uri"`
}

var (
	onetext, hitokoto, april,  all       			onetextData
	i1, i2, i3, i, randNub                  		int
	jsonByte, jsonData1, jsonData2, jsonData3		[]byte
	err                             				error
	jsonStr											string
)

func dropErr(e error) {
	if e != nil {
		panic(e)
	}
}

// FormatJSON Fucking VS Code let me add a comment
func FormatJSON(Source string) {
	switch {
	case Source == "onetext":
		all = onetext
		i = i1
	case Source == "hitokoto":
		all = hitokoto
		i = i2
	case Source == "april":
		all = april
		i = i3
	default:
		all = onetext
		i = i1
	}
	rand.Seed(time.Now().UnixNano())
	randNub := rand.Intn(i)
	if all[randNub].URI != "" {
		jsonRaw := struct {
			Text string   `json:"text"`
			By   string   `json:"by"`
			From string   `json:"from"`
			Time []string `json:"time"`
			URI  string   `json:"uri"`
		}{all[randNub].Text, all[randNub].By, all[randNub].From, all[randNub].Time, all[randNub].URI}
		jsonByte ,err := json.MarshalIndent(jsonRaw, "", "    ")
		if err != nil {
			fmt.Println(err)
		}
		jsonStr = string(jsonByte)
		return
	}
	jsonRaw := struct {
		Text string   `json:"text"`
		By   string   `json:"by"`
		From string   `json:"from"`
		Time []string `json:"time"`
	}{all[randNub].Text, all[randNub].By, all[randNub].From, all[randNub].Time}
	jsonByte ,err := json.MarshalIndent(jsonRaw, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	jsonStr = string(jsonByte)
	return
}

// ResponseOnetext Fucking VS Code let me add a comment
func ResponseOnetext(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	r.ParseForm()
	fmt.Println(r.Form)
	var Source string
	if _, ok := r.Form["source"];ok {
		Source = string(r.Form["source"][0])
	}
	FormatJSON(Source)
	fmt.Println(string(jsonStr))
	fmt.Fprintf(w, strings.Replace(string(jsonStr), "%", "%%", -1))
}

func main() {

	jsonFile1 := "./OneText-Library.json"
	jsonFile2 := "./all.json"
	jsonFile3 := "./april.json"

	jsonData1, err := ioutil.ReadFile(jsonFile1)
	dropErr(err)
	jsonData2, err := ioutil.ReadFile(jsonFile2)
	dropErr(err)
	jsonData3, err := ioutil.ReadFile(jsonFile3)
	dropErr(err)

	i1 = strings.Count(string(jsonData1), "\"text\"")
	i2 = strings.Count(string(jsonData2), "\"text\"")
	i3 = strings.Count(string(jsonData3), "\"text\"")

	err = json.Unmarshal([]byte(jsonData1), &onetext)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal([]byte(jsonData2), &hitokoto)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal([]byte(jsonData3), &april)
	if err != nil {
		fmt.Println(err)
	}

	http.HandleFunc("/v1", ResponseOnetext)
	http.HandleFunc("/api", ResponseOnetext)
	http.ListenAndServe(":8000", nil)

}
