package main

import (

	//CALCULATE THE SHORTEST DISTANCE
	// "aqiCrawler/distance"
	// "fmt"

	"aqiCrawler/db"
	"aqiCrawler/linebot"
	"fmt"
	"net/http"
	"os"
)

func main() {

	//GET JSON INSERT DB
	aqidb.GetData()

	//CALCULATE THE SHORTEST DISTANCE
	// str := distance.GetSite("25.100000", "121.500000")
	// fmt.Println(str)

	app, err := linebot.NewLineBot(
		os.Getenv("ChannelSecret"),
		os.Getenv("ChannelAccessToken"),
		os.Getenv("APP_BASE_URL"),
	)
	if err != nil {
		fmt.Println(err)
	}
	http.HandleFunc("/callback", app.Callback)
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		fmt.Println(err)
	}
}
