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
	http.HandleFunc("/test_callback", app.Callback)
	http.HandleFunc("/test_auth", linenotify.Auth)
	http.HandleFunc("/test_pushnotify", linenotify.Token)
	http.HandleFunc("/test_getdata", db.GetData)
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		fmt.Println(err)
	}
}
