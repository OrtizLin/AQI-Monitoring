package db

import (
	"aqiCrawler/linenotify"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

type User struct {
	UserClientID string
	UserToken    string
	UserLocation []string
}

func GetData() {
	//Connect DB
	session, errs := mgo.Dial(os.Getenv("DBURL"))
	if errs != nil {
		panic(errs)
	}
	defer session.Close()
	c := session.DB("aqidb").C("aqisite")
	c2 := session.DB("aqidb").C("userdb")
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
		//Check status and send notify to whom live in this area.
		if aqisite.Status == "良好" {
			result := User{}
			iter := c2.Find(nil).Iter()
			for iter.Next(&result) {
				if contains(result.UserLocation, aqisite.StieName) {
					linenotify.SendGoodNews(result.UserToken, aqisite.StieName)
				}
			}
		}
	}
}

func SaveToken(token, clientID string) bool {
	session, errs := mgo.Dial(os.Getenv("DBURL"))
	if errs != nil {
		panic(errs)
	}
	defer session.Close()
	collect := session.DB("aqidb").C("userdb")
	loc := []string{}
	user := User{}
	user.UserToken = token
	user.UserClientID = clientID
	user.UserLocation = loc
	errs = collect.Insert(&User{user.UserClientID, user.UserToken, user.UserLocation})
	if errs != nil {
		return false
	} else {
		return true
	}
}

func NewSite(site, clientID string) bool {
	result := User{}
	session, errs := mgo.Dial(os.Getenv("DBURL"))
	if errs != nil {
		panic(errs)
	}
	defer session.Close()
	c := session.DB("aqidb").C("userdb")

	err := c.Find(bson.M{"userclientid": clientID}).One(&result)
	if err != nil {
		fmt.Println("NOT FOUND!!")
		return false
	} else {
		fmt.Println("FOUND!!!")
		if contains(result.UserLocation, site) {
			return true
		} else {
			result.UserLocation = append(result.UserLocation, site)
			colQuerier := bson.M{"userclientid": clientID}
			change := bson.M{"$set": bson.M{"userlocation": result.UserLocation}}
			err = c.Update(colQuerier, change)
			if err != nil {
				panic(err)
			}
			return true
		}
	}
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
