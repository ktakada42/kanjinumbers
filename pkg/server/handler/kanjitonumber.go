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

// var separatorsEveryFourDigit = []int{1000000000000, 100000000, 10000, 1} // 1兆, 1億, 1万 (numberstokanji.goで定義済)
// var separatorsOfFourDigit = []int{1000, 100, 10, 1} // numberstokanji.goで定義済

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

	num, err := ConvertKanjiToNumber(arr[3])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	fmt.Fprintf(w, strconv.Itoa(num))
}

func ConvertKanjiToNumber(kanji string) (int, error) {
	var num int
	var err error

	if kanji == "零" {
		return num, err
	}

	// 有効なパラメーターの最大長は31文字
	if len([]rune(kanji)) > 31 {
		err = errors.New("Param invalid")
		return num, err
	}

	/* ===== パラメーターを上から4桁ずつ配列に入れる ===== */
	// Ex) "壱千弐百参拾四億五千六百七拾八" => {"", "壱千弐百参拾四", "", "五千六百七拾八"}
	kanjiSeparatedEveryFourDigit, err := separateKanjiEveryFourDigit(kanji)
	if err != nil {
		return num, err
	}

	for i, s := range kanjiSeparatedEveryFourDigit {

		/* ===== 4桁ずつ分けた数字を更に1桁ずつに区切る ===== */
		kanjiSeparatedEveryDigit, err := separateKanjiEveryDigit(s)
		if err != nil {
			return num, err
		}
		for j, s := range kanjiSeparatedEveryDigit {

			/* ===== 1桁ずつに区切った漢数字を数字に変換 ===== */
			n, err := convertEveryKanjiToNumber(s)
			if err != nil {
				return num, err
			}

			// それぞれの数字に区切り桁を掛け、戻り値に足す
			num += n * separatorsOfFourDigit[j] * separatorsEveryFourDigit[i]
		}
	}
	return num, err
}

// パラメーターを上から4桁ずつ配列に入れる
// Ex) 壱千弐百参拾四億五千六百七拾八 => {"", "壱千弐百参拾四", "", "五千六百七拾八"}

// この時点では、桁区切り文字の前に何らかの文字列があるかどうかのみをValidate
// Ex) "兆", "壱兆億", "兆壱億" => 204
func separateKanjiEveryFourDigit(kanji string) (kanjiSeparatedEveryFourDigit [4]string, err error) {
	arr := strings.SplitN(kanji, "兆", 2)
	if len(arr) == 2 {
		if arr[0] == "" {
			err = errors.New("Param invalid")
			return
		}
		kanjiSeparatedEveryFourDigit[0] = arr[0]
		kanji = arr[1]
	}
	arr = strings.SplitN(kanji, "億", 2)
	if len(arr) == 2 {
		if arr[0] == "" {
			err = errors.New("Param invalid")
			return
		}
		kanjiSeparatedEveryFourDigit[1] = arr[0]
		kanji = arr[1]
	}
	arr = strings.SplitN(kanji, "万", 2)
	if len(arr) == 2 {
		if arr[0] == "" {
			err = errors.New("Param invalid")
			return
		}
		kanjiSeparatedEveryFourDigit[2] = arr[0]
		kanjiSeparatedEveryFourDigit[3] = arr[1]
		return
	}
	kanjiSeparatedEveryFourDigit[3] = arr[0]
	return
}

// 4桁ずつ分けた数字を更に1桁ずつに区切る

// この時点では、桁区切り文字の前に何らかの文字列があるかどうかのみをValidate
// Ex) "千", "壱千百", "千壱百" => 204
func separateKanjiEveryDigit(kanjiSeparatedEveryFourDigit string) (kanjiSeparatedEveryDigit [4]string, err error) {
	arr := strings.SplitN(kanjiSeparatedEveryFourDigit, "千", 2)
	if len(arr) == 2 {
		if arr[0] == "" {
			err = errors.New("Param invalid")
			return
		}
		kanjiSeparatedEveryDigit[0] = arr[0]
		kanjiSeparatedEveryFourDigit = arr[1]
	}
	arr = strings.SplitN(kanjiSeparatedEveryFourDigit, "百", 2)
	if len(arr) == 2 {
		if arr[0] == "" {
			err = errors.New("Param invalid")
			return
		}
		kanjiSeparatedEveryDigit[1] = arr[0]
		kanjiSeparatedEveryFourDigit = arr[1]
	}
	arr = strings.SplitN(kanjiSeparatedEveryFourDigit, "拾", 2)
	if len(arr) == 2 {
		if arr[0] == "" {
			err = errors.New("Param invalid")
			return
		}
		kanjiSeparatedEveryDigit[2] = arr[0]
		kanjiSeparatedEveryDigit[3] = arr[1]
		return
	}
	kanjiSeparatedEveryDigit[3] = arr[0]
	return
}

// 1桁ずつに区切った漢数字を数字に変換
func convertEveryKanjiToNumber(digit string) (num int, err error) {

	// 1桁ずつに区切った文字列がそれぞれ1文字 or 空文字かをValidate
	// Ex) {"", "壱", "", "弐"} => OK
	// Ex) {"壱弐", "", "", ""} => 204
	digitLen := len([]rune(digit))
	if digitLen > 1 {
		err = errors.New("Param invalid")
		return num, err
	}

	// 1文字の文字列に関し、有効な漢数字1文字であるかをVaildate
	// Ex) {"一", "", "", ""} => 204
	if digitLen == 1 {
		num, isValidChar := numKanji[digit]
		if !isValidChar {
			err = errors.New("Param invalid")
		}
		return num, err
	}
	return num, err
}
