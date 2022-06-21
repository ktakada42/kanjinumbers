package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func HandleKanjiToNumber(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.Path)
	if err != nil {
		log.Println("Parse failed")
		return
	}
	arr := strings.Split(u.Path, "/")
	if len(arr) != 4 {
		log.Println("Path invalid")
		return
	}
	if !isPathParamValid(arr[3]) {
		log.Println("Param invalid")
		return
	}
	fmt.Fprintf(w, arr[3])
}
