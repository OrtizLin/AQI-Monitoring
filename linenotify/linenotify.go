package linenotify

import (
	"fmt"
	"github.com/utahta/go-linenotify"
	"github.com/utahta/go-linenotify/auth"
	//"github.com/utahta/go-linenotify/token"
	"net/http"
	"os"
	"time"
)

func Auth(w http.ResponseWriter, req *http.Request) {
	param1 := req.URL.Query().Get("client")
	fmt.Fprint(w, param1)

	c, err := auth.New(os.Getenv("ClientID"), os.Getenv("APP_BASE_URL")+"pushnotify?client="+param1)
	if err != nil {
		fmt.Fprintf(w, "error:%v", err)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "state", Value: c.State, Expires: time.Now().Add(60 * time.Second)})
	c.Redirect(w, req)
}

func Token(w http.ResponseWriter, req *http.Request) {
	param1 := req.URL.Query().Get("client")
	fmt.Fprint(w, param1)
}