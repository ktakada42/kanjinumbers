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

var kanjiNumbers = map[string]string{"0": "零", "1": "壱", "2": "弐", "3": "参", "4": "四", "5": "五", "6": "六", "7": "七", "8": "八", "9": "九"}

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
	fmt.Fprintf(w, convertNumberToKanji(arr[3]))
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

func convertNumberToKanji(num string) (kanji string) {
	n, _ := strconv.Atoi(num)
	if n == 0 {
		kanji = "零"
		return
	}
	digit := len(num)
	revNum := strRev(num)
	for digitCount, n := range revNum {
		if n != '0' {
			switch digitCount % 4 {
			case 1:
				kanji = "拾" + kanji
			case 2:
				kanji = "百" + kanji
			case 3:
				kanji = "千" + kanji
			default:
				if digitCount != digit {
					switch digitCount / 4 {
					case 1:
						kanji = "万" + kanji
					case 2:
						kanji = "億" + kanji
					case 3:
						kanji = "兆" + kanji
					default:
					}
				}
			}
			kanji = kanjiNumbers[string(n)] + kanji
		}
	}
	return
}
