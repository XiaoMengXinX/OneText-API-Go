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
	onetext, hitokoto, april	       			onetextData
	i1, i2, i3, i, randNub                  		int
	jsonByte, jsonData1, jsonData2, jsonData3		[]byte
	err                             				error
	jsonStr, textStr ,ResponseData				string
)

func dropErr(e error) {
	if e != nil {
		panic(e)
	}
}

// FormatJSON Fucking VS Code let me add a comment
func FormatJSON(source onetextData, i int) {
	rand.Seed(time.Now().UnixNano())
	randNub := rand.Intn(i)
	if source[randNub].URI != "" {
		jsonRaw := struct {
			Text string   `json:"text"`
			By   string   `json:"by"`
			From string   `json:"from"`
			Time []string `json:"time"`
			URI  string   `json:"uri"`
		}{source[randNub].Text, source[randNub].By, source[randNub].From, source[randNub].Time, source[randNub].URI}
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
	}{source[randNub].Text, source[randNub].By, source[randNub].From, source[randNub].Time}
	jsonByte ,err := json.MarshalIndent(jsonRaw, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	jsonStr = string(jsonByte)
	return
}

// GetText Fucking VS Code let me add a comment
func GetText(source onetextData, i int) {
	rand.Seed(time.Now().UnixNano())
	randNub := rand.Intn(i)
	textStr = source[randNub].Text
	return
}

// ResponseOnetext Fucking VS Code let me add a comment
func ResponseOnetext(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r)
	fmt.Println(r.Form)
	var Source ,dataType string
	var source onetextData
	if _, ok := r.Form["source"];ok {
		Source = string(r.Form["source"][0])
	}
	if _, ok := r.Form["type"];ok {
		dataType = string(r.Form["type"][0])
	}
	switch {
	case Source == "onetext":
		source = onetext
		i = i1
	case Source == "hitokoto":
		source = hitokoto
		i = i2
	case Source == "april":
		source = april
		i = i3
	default:
		source = onetext
		i = i1
	}
	switch {
	case dataType == "json":
		FormatJSON(source, i)
		ResponseData = string(jsonStr)
		w.Header().Set("content-type", "application/json; charset=utf-8")
	case dataType == "text":
		GetText(source, i)
		ResponseData = string(textStr)
		w.Header().Set("content-type", "application/json; charset=utf-8")
	default:
		FormatJSON(source, i)
		ResponseData = string(jsonStr)
		w.Header().Set("content-type", "text/plain; charset=utf-8")
	}
	w.WriteHeader(200)
	fmt.Println(ResponseData)
	fmt.Fprintf(w, strings.Replace(ResponseData, "%", "%%", -1))
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
