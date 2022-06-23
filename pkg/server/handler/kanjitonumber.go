package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var numKanji = map[string]int{"壱": 1, "弐": 2, "参": 3, "四": 4, "五": 5, "六": 6, "七": 7, "八": 8, "九": 9}

// var largeSeparatorsNum = []int{1000000000000, 100000000, 10000, 1} // 1兆, 1億, 1万(numberstokanji.goで定義済)
var smallSeparatorsNum = []int{1000, 100, 10, 1}

func HandleKanjiToNumber(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.Path)
	if err != nil {
		log.Println("Parse failed")
		return
	}
	arr := strings.Split(u.Path, "/")

	// パラメーターの数が2つ以上か、パラメーターが無いと204
	// Ex) /v1/kanji2number/壱百弐拾参/四百五拾六 => 204
	// Ex) /v1/kanji2number/ => 204
	if len(arr) != 4 || arr[3] == "" {
		log.Println("Path invalid")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	num, err := convertKanjiToNumber(arr[3])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	fmt.Fprintf(w, strconv.Itoa(num))
}

func convertKanjiToNumber(kanji string) (num int, err error) {
	if kanji == "零" {
		return
	}

	// 有効なパラメーターの最大長は31文字
	if len([]rune(kanji)) > 31 {
		err = errors.New("Param invalid")
		return
	}

	// パラメーターを上から4桁ずつ配列に入れる
	// Ex) 壱千弐百参拾四億五千六百七拾八 => {"", "壱千弐百参拾四", "", "五千六百七拾八"}

	// この時点では、桁区切り文字の前に何らかの文字列があるかどうかのみをValidate
	// Ex) "兆", "壱兆億", "兆壱億" => 204
	var largeSeparatedKanji [4]string
	arr := strings.SplitN(kanji, "兆", 2)
	if len(arr) == 2 {
		if arr[0] == "" {
			err = errors.New("Param invalid")
			return
		}
		largeSeparatedKanji[0] = arr[0]
		kanji = arr[1]
	}
	arr = strings.SplitN(kanji, "億", 2)
	if len(arr) == 2 {
		if arr[0] == "" {
			err = errors.New("Param invalid")
			return
		}
		largeSeparatedKanji[1] = arr[0]
		kanji = arr[1]
	}
	arr = strings.SplitN(kanji, "万", 2)
	if len(arr) == 2 {
		if arr[0] == "" {
			err = errors.New("Param invalid")
			return
		}
		largeSeparatedKanji[2] = arr[0]
		largeSeparatedKanji[3] = arr[1]
	} else {
		largeSeparatedKanji[3] = arr[0]
	}

	// 4桁ずつ分けた数字を更に1桁ずつに区切る

	// この時点では、桁区切り文字の前に何らかの文字列があるかどうかのみをValidate
	// Ex) "千", "壱千百", "千壱百" => 204
	for i, s := range largeSeparatedKanji {
		var smallSeparatedKanji [4]string
		arr := strings.SplitN(s, "千", 2)
		if len(arr) == 2 {
			if arr[0] == "" {
				err = errors.New("Param invalid")
				return
			}
			smallSeparatedKanji[0] = arr[0]
			s = arr[1]
		}
		arr = strings.SplitN(s, "百", 2)
		if len(arr) == 2 {
			if arr[0] == "" {
				err = errors.New("Param invalid")
				return
			}
			smallSeparatedKanji[1] = arr[0]
			s = arr[1]
		}
		arr = strings.SplitN(s, "拾", 2)
		if len(arr) == 2 {
			if arr[0] == "" {
				err = errors.New("Param invalid")
				return
			}
			smallSeparatedKanji[2] = arr[0]
			smallSeparatedKanji[3] = arr[1]
		} else {
			smallSeparatedKanji[3] = arr[0]
		}

		for j, s := range smallSeparatedKanji {
			_, isValidChar := numKanji[s]

			// 1桁ずつに区切った文字列がそれぞれ空文字か、或いは一文字かつ有効な漢数字であるかをVaildate
			// Ex) {"", "壱", "", "弐"} => OK
			// Ex) {"壱弐", "", "", ""} => 204
			// Ex) {"一", "", "", ""} => 204
			if len([]rune(s)) > 1 || s != "" && !isValidChar {
				err = errors.New("Param invalid")
				return
			}

			// 漢数字を対応する数字に変換し、区切り桁を掛けて戻り値に足す
			if s != "" {
				n := numKanji[s]
				num += (n * smallSeparatorsNum[j] * largeSeparatorsNum[i])
			}
		}
	}
	return
}
