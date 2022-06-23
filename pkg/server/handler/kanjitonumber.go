package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var validChar string = "壱弐参四五六七八九拾百千万億兆"

func HandleKanjiToNumber(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.Path)
	if err != nil {
		log.Println("Parse failed")
		return
	}
	arr := strings.Split(u.Path, "/")

	// パラメーターの数が合わないと204
	// Ex) /v1/kanji2number/壱百二拾三/四百五拾六 => 204
	// Ex) /v1/kanji2number/ => 204
	if len(arr) != 4 || arr[3] == "" {
		log.Println("Path invalid")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	fmt.Fprintf(w, convertKanjiToNumber(arr[3]))
}

func convertKanjiToNumber(kanji string) (num string) {
	if kanji == "零" {
		num = "0"
		return
	}
	if isContainedOnlyValidChar(kanji) {
		return "true"
	}
	return "false"
}

func isContainedOnlyValidChar(kanji string) bool {
	haystack := []rune(validChar)
	needles := []rune(kanji)
	for _, n := range needles {
		found := false
		for _, h := range haystack {
			if h == n {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}
