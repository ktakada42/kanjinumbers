package handler

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	min = 0
	max = 9999999999999999
)

var kanjiNumbers = map[int]string{0: "", 1: "壱", 2: "弐", 3: "参", 4: "四", 5: "五", 6: "六", 7: "七", 8: "八", 9: "九"}
var largeSeparatorsNum = []int{1000000000000, 100000000, 10000, 1} // 1兆, 1億, 1万
var largeSeparatorsKanji = []string{"兆", "億", "万", ""}             // 桁区切り(大)
var smallSeparators = []string{"", "拾", "百", "千"}                  // 桁区切り(小)

func HandleNumberToKanji(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.Path)
	if err != nil {
		log.Println("Parse failed")
		return
	}
	arr := strings.Split(u.Path, "/")

	// パラメーターの数が合わないと204
	// Ex) /v1/number2kanji/123/123 => 204
	// Ex) /v1/number2kanji/ => 204
	if len(arr) != 4 {
		log.Println("Path invalid")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// 変換できないパラメーターは204
	// Ex) /v1/number2kanji/a => 204
	// Ex) /v1/number2kanji/-1 => 204
	// Ex) /v1/number2kanji/9999999999999999999999999999 => 204
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
	// '+'はURIの予約文字なので弾く
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

	// パラメーターを上から4桁ずつ配列に入れる
	// Ex) 123,456,789 => [0(兆), 1(億), 2345(万), 6789]
	var separatedNum [4]int
	for i, s := range largeSeparatorsNum {
		separatedNum[i] = n / s
		n %= s
	}

	// 4桁ずつ分けた数字を更に桁区切り(小)で区切る
	for i, s := range separatedNum {
		if s != 0 {
			for i := 3; i >= 0; i-- {
				if s/int(math.Pow10(i)) != 0 {
					kanji += kanjiNumbers[s/int(math.Pow10(i))] + smallSeparators[i]
				}
				s %= int(math.Pow10(i))
			}
			kanji += largeSeparatorsKanji[i] // 桁区切り(大)を追加
		}
	}
	return
}
