package main

import (
	"aqiCrawler/db"
	"aqiCrawler/linebot"
	"aqiCrawler/linenotify"
	"fmt"
	"net/http"
	"os"
)

func main() {
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
	http.HandleFunc("/pushnotify", linenotify.Token)
	http.HandleFunc("/getdata", db.GetData)
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		fmt.Println(err)
	}
}
