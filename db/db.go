package db

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os"
)

type AqiSite struct {
	StieName   string
	AQI        string
	Status     string
	Latitude   string
	Longitude  string
	UpdateTime string
}

func GetData() {
	//Connect DB
	session, errs := mgo.Dial(os.Getenv("DBURL"))
	if errs != nil {
		panic(errs)
	}
	defer session.Close()
	c := session.DB("aqidb").C("aqisite")
	//Clean DB
	c.RemoveAll(nil)
	//Get AQI data from opendate2
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
		aqisite := AqiSite{}
		md := generic[i].(map[string]interface{})
		aqisite.StieName = md["SiteName"].(string)
		aqisite.AQI = md["AQI"].(string)
		aqisite.Status = md["Status"].(string)
		aqisite.Latitude = md["Latitude"].(string)
		aqisite.Longitude = md["Longitude"].(string)
		aqisite.UpdateTime = md["PublishTime"].(string)
		//Insert to DB
		errs = c.Insert(&AqiSite{aqisite.StieName, aqisite.AQI, aqisite.Status, aqisite.Latitude, aqisite.Longitude, aqisite.UpdateTime})
		if errs != nil {
			log.Fatal(errs)
		}
	}
}
