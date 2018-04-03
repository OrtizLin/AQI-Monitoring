package linenotify

import (
	"fmt"
	"net/http"
)

func Auth(w http.ResponseWriter, req *http.Request) {
	param1 := req.URL.Query().Get("client")
	fmt.Fprint(w, param1)
}
