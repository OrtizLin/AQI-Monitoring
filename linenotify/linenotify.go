package linenotify

import (
	"aqiCrawler/db"
	"fmt"
	"github.com/utahta/go-linenotify"
	"github.com/utahta/go-linenotify/auth"
	"github.com/utahta/go-linenotify/token"
	"net/http"
	"os"
	"time"
)

type User struct {
	UserClientID string
	UserToken    string
	UserLocation []string
}

func Auth(w http.ResponseWriter, req *http.Request) {
	param1 := req.URL.Query().Get("client")
	c, err := auth.New(os.Getenv("ClientID"), os.Getenv("APP_BASE_URL")+"pushnotify")
	if err != nil {
		fmt.Fprintf(w, "error:%v", err)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "state", Value: param1, Expires: time.Now().Add(60 * time.Second)})
	c.Redirect(w, req)
}

func Token(w http.ResponseWriter, req *http.Request) {
	resp, err := auth.ParseRequest(req)
	if err != nil {
		fmt.Fprintf(w, "error:%v", err)
		return
	}

	state, err := req.Cookie("state")
	if err != nil {
		fmt.Fprintf(w, "error:%v", err)
		return
	}
	c := token.New(os.Getenv("APP_BASE_URL")+"pushnotify", os.Getenv("ClientID"), os.Getenv("ClientSecret"))
	accessToken, err := c.GetAccessToken(resp.Code)
	if err != nil {
		fmt.Fprintf(w, "error:%v", err)
		return
	}

	if db.SaveToken(accessToken, state.Value) {
		fmt.Fprintf(w, "LINE Notify 連動完成。\n 請返回空汙報報, 並註冊離你最近的觀測站。")
	}

}

func SomeOneFollow(displayname, url string) {
	token := os.Getenv("OtisToken")
	c := linenotify.New()
	content := displayname + "追蹤了空汙報報 ！"
	c.NotifyWithImageURL(token, content, url, url)
}
