package distance

import (
	"gopkg.in/mgo.v2"
	"math"
	"os"
	"strconv"
)

type AqiSite struct {
	SiteName   string
	AQI        string
	Status     string
	Latitude   string
	Longitude  string
	UpdateTime string
}

func GetSite(lat, lon string) string {
	dis, err := strconv.ParseFloat("10000000", 64)
	site := ""
	lat2, err := strconv.ParseFloat(lat, 64) //User's Lat
	lon2, err := strconv.ParseFloat(lon, 64) //User's Long
	if err != nil {
		panic(err)
	}
	session, errs := mgo.Dial(os.Getenv("DBURL"))
	if errs != nil {
		panic(errs)
	}
	defer session.Close()
	c := session.DB("aqidb").C("aqisite")

	result := AqiSite{}
	iter := c.Find(nil).Iter()
	for iter.Next(&result) {

		lat1, err := strconv.ParseFloat(result.Latitude, 64)
		lon1, err := strconv.ParseFloat(result.Longitude, 64)
		if err != nil {
			panic(err)
		}
		rangeOfTwoPoint := distance(lat1, lon1, lat2, lon2)
		if rangeOfTwoPoint < dis {
			dis = rangeOfTwoPoint
			site = result.SiteName
		}
	}
	return site
}

func distance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}
