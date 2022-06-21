package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	min = 0
	max = 9999999999999999
)

func HandleNumberToKanji(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.Path)
	if err != nil {
		log.Println("Parse failed")
		return
	}
	arr := strings.Split(u.Path, "/")
	if len(arr) != 4 {
		log.Println("Path invalid")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if !isPathParamValid(arr[3]) {
		log.Println("Param invalid")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	fmt.Fprintf(w, arr[3])
}

func isPathParamValid(p string) bool {
	n, err := strconv.Atoi(p)
	if err != nil {
		return false
	}
	if p[0] == '+' {
		return false
	}
	if min <= n && n <= max {
		return true
	}
	return false
}
