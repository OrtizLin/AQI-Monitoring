package main

import (
	"apiCrawler/linenotify"
	"aqiCrawler/db"
	"aqiCrawler/linebot"
	"fmt"
	"net/http"
	"os"
)

func main() {

	//GET JSON INSERT DB
	db.GetData()

	app, err := linebot.NewLineBot(
		os.Getenv("ChannelSecret"),
		os.Getenv("ChannelAccessToken"),
		os.Getenv("APP_BASE_URL"),
	)
	if err != nil {
		fmt.Println(err)
	}
	http.HandleFunc("/callback", app.Callback)
	http.HandleFunc("/auth", linenotify.Auth)
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		fmt.Println(err)
	}
}
