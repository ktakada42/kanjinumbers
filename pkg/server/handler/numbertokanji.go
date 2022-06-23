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
	Min = 0
	Max = 9999999999999999
)

var kanjiNum = map[int]string{1: "壱", 2: "弐", 3: "参", 4: "四", 5: "五", 6: "六", 7: "七", 8: "八", 9: "九"}
var separatorsEveryFourDigit = []int{1000000000000, 100000000, 10000, 1} // 1兆, 1億, 1万
var separatorsOfFourDigit = []int{1000, 100, 10, 1}
var kanjiSeparatorsEveryFourDigit = []string{"兆", "億", "万", ""}
var kanjiSeparatorsOfFourDigit = []string{"千", "百", "拾", ""}

func HandleNumberToKanji(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.Path)
	if err != nil {
		log.Println("Parse failed")
		return
	}
	arr := strings.Split(u.Path, "/")

	// パラメーターの数が2つ以上か、パラメーターが無いと204
	// Ex) /v1/number2kanji/123/456 => 204
	// Ex) /v1/number2kanji/ => 204
	if len(arr) != 4 || arr[3] == "" {
		log.Println("Path invalid")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	/* ===== パラメーターをValidate ===== */
	// Ex) /v1/number2kanji/a => 204
	// Ex) /v1/number2kanji/-1 => 204
	// Ex) /v1/number2kanji/9999999999999999999999999999 => 204
	if !isPathParamValid(arr[3]) {
		log.Println("Param invalid")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	fmt.Fprintf(w, ConvertNumberToKanji(arr[3]))
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
	if Min <= n && n <= Max {
		return true
	}
	return false
}

func ConvertNumberToKanji(num string) (kanji string) {
	n, _ := strconv.Atoi(num)
	if n == 0 {
		kanji = "零"
		return
	}

	// パラメーターを上から4桁ずつ配列に入れる
	// Ex) 123,456,789 => {"0(兆)", "1(億)", "2345(万)", "6789"}
	var numSeparatedEveryFourDigit [4]int
	for i, s := range separatorsEveryFourDigit {
		numSeparatedEveryFourDigit[i] = n / s
		n %= s
	}

	// 4桁ずつ分けた数字を更に1桁ずつに区切って漢数字に変換、区切り桁を足す
	// Ex) "2345(万)" => "弐" + "千" + "参" + "百" + "四" + "拾" + "五" + "万"
	for i, s := range numSeparatedEveryFourDigit {
		if s != 0 {
			for i := 0; i < 4; i++ {
				sep := separatorsOfFourDigit[i]
				if s/sep != 0 {
					kanji += kanjiNum[s/sep] + kanjiSeparatorsOfFourDigit[i] // 漢数字に変換 + 「千 or 百 or 拾」を追加
				}
				s %= sep
			}
			kanji += kanjiSeparatorsEveryFourDigit[i] // 区切り文字として、「兆 or 億 or 万」を追加
		}
	}
	return
}
