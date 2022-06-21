package handler

import (
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	min       = 0
	max       = 9999999999999999
	kanji_num = "零壱弐参四五六七八九"
	sep       = "零万億兆"
)

// HandleNumberToKanji 数字から漢数字への変換
func HandleNumberToKanji() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}

func numberToKanji(req *http.Request) string {
	u, err := url.Parse(req.URL.Path)
	if err != nil {
		return "Parse failed"
	}
	arr := strings.Split(u.Path, "/")
	if len(arr) != 4 {
		return "Path invalid"
	}
	if !isPathParamValid(arr[3]) {
		return "Param invalid"
	}
	return transferNumberToKanji(arr[3])
}

func isPathParamValid(p string) bool {
	n, err := strconv.Atoi(p)
	if err != nil {
		return false
	}
	if min <= n && n <= max {
		return true
	}
	return false
}

func transferNumberToKanji(num string) (kanji string) {
	n, _ := strconv.Atoi(num)
	if n == 0 {
		return "零"
	}
	for i := 12; i >= 4; i -= 4 {
		n /= int(math.Pow10(i))
	}

	// n := strconv.Atoi(num)
	// for i := 15; i > 0; i++ {
	// 	if n/math.Pow(10, i) >= 1 {

	// 	}
	// }
}

func splitEveryFourDigits(num string) (sepNum []string) {
	for i := len(num) - 1; i > 0; i-- {

	}
}
