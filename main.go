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
	http.HandleFunc("/remove_callback", app.Callback)
	http.HandleFunc("/remove_auth", linenotify.Auth)
	http.HandleFunc("/remove_pushnotify", linenotify.Token)
	http.HandleFunc("/remove_getdata", db.GetData)
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		fmt.Println(err)
	}
}
