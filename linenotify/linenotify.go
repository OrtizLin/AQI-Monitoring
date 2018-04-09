package linenotify

import (
	"fmt"
	//"github.com/utahta/go-linenotify"
	"github.com/utahta/go-linenotify/auth"
	"github.com/utahta/go-linenotify/token"
	"net/http"
	"os"
	"time"
)

const param1 string = ""

func Auth(w http.ResponseWriter, req *http.Request) {
	param1 = req.URL.Query().Get("client")
	c, err := auth.New(os.Getenv("ClientID"), os.Getenv("APP_BASE_URL")+"pushnotify")
	if err != nil {
		fmt.Fprintf(w, "error:%v", err)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "state", Value: c.State, Expires: time.Now().Add(60 * time.Second)})
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
	if resp.State != state.Value {
		fmt.Fprintf(w, "error:%v", err)
		return
	}

	c := token.New(os.Getenv("APP_BASE_URL")+"pushnotify", os.Getenv("ClientID"), os.Getenv("ClientSecret"))
	accessToken, err := c.GetAccessToken(resp.Code)
	if err != nil {
		fmt.Fprintf(w, "error:%v", err)
		return
	}
	fmt.Println(accessToken, "client id is :", param1)
}
