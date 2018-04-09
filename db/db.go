package db

import (
	"encoding/json"
	"fmt"
	"github.com/utahta/go-linenotify"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
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

func GetData(w http.ResponseWriter, req *http.Request) {
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

		time := aqisite.UpdateTime[len(aqisite.UpdateTime)-5:]
		//Only notify at 8 am and 6 pm
		if time == "08:00" || time == "18:00" {
			//Check status and send notify to whom live in this area.
			if aqisite.Status == "良好" {
				result := User{}
				iter := c2.Find(nil).Iter()
				for iter.Next(&result) {
					if contains(result.UserLocation, aqisite.StieName) {
						connect := linenotify.New()
						//Random pokémon pic
						myrand := random(1, 251)
						url := ""
						if myrand < 10 {
							url = "https://www.dragonflycave.com/sprites/gen2/g/00" + strconv.Itoa(myrand) + ".png"
						} else if myrand >= 10 && myrand < 100 {
							url = "https://www.dragonflycave.com/sprites/gen2/g/0" + strconv.Itoa(myrand) + ".png"
						} else {
							url = "https://www.dragonflycave.com/sprites/gen2/g/" + strconv.Itoa(myrand) + ".png"
						}
						str := "今天 [" + aqisite.StieName + "] 附近空氣良好, 把握機會出去走走吧！"
						connect.NotifyWithImageURL(result.UserToken, str, url, url)
					}
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

func CheckRegistered(clientID string) string {
	result := User{}
	session, errs := mgo.Dial(os.Getenv("DBURL"))
	if errs != nil {
		panic(errs)
	}
	defer session.Close()
	c := session.DB("aqidb").C("userdb")
	err := c.Find(bson.M{"userclientid": clientID}).One(&result)
	if err != nil {
		return "NoUser"
	} else {
		str := ""
		for i := 0; i < len(result.UserLocation); i++ {
			if i == len(result.UserLocation)-1 {
				str = str + result.UserLocation[i]
			} else {
				str = str + result.UserLocation[i] + ","
			}
		}
		return str
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

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
