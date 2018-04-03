package linenotify

import (
	"fmt"
	"net/http"
)

func Auth(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "hello world ")
}
