package main

import (
	"aqiCrawler/distancepackage"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	dis, err := strconv.ParseFloat("10000000", 64)
	if err != nil {
	}
	site := ""

	resp, err := http.Get("http://opendata2.epa.gov.tw/AQI.json")
	if err != nil {
		log.Fatal(err)
	}
	var generic []interface{}
	err = json.NewDecoder(resp.Body).Decode(&generic)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(generic); i++ {
		md := generic[i].(map[string]interface{})
		lat1, err := strconv.ParseFloat(md["Latitude"].(string), 64)
		lon1, err := strconv.ParseFloat(md["Longitude"].(string), 64)
		lat2, err := strconv.ParseFloat("25.048220", 64)  //User's Lat
		lon2, err := strconv.ParseFloat("121.532731", 64) //User's Long
		if err != nil {

		}
		fmt.Println(md["SiteName"], md["AQI"], md["Status"], md["Latitude"], md["Longitude"])
		rangeOfTwoPoint := distancepackage.Distance(lat1, lon1, lat2, lon2)
		if rangeOfTwoPoint < dis {
			dis = rangeOfTwoPoint
			site = md["SiteName"].(string)
		}
	}
	fmt.Println("The shortest site is :", site, "range is :", dis)
}
