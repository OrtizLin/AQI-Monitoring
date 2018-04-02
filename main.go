package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	resp, err := http.Get("http://opendata2.epa.gov.tw/AQI.json")
	if err != nil {
		log.Fatal(err)
	}
	var generic []interface{}
	err = json.NewDecoder(resp.Body).Decode(&generic)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 77; i++ {
		md := generic[i].(map[string]interface{})
		fmt.Println(md["SiteName"], md["AQI"], md["Status"], md["Latitude"], md["Longitude"])
	}
}
