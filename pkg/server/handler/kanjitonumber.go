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
var smallSeparatorsNum = []int{1000, 100, 10, 1}

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
	if len([]rune(kanji)) > 31 {
		err = errors.New("Param invalid")
		return
	}
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
			if len([]rune(s)) > 1 || s != "" && !isValidChar {
				err = errors.New("Param invalid")
				return
			}
			if s != "" {
				n := numKanji[s]
				num += (n * smallSeparatorsNum[j] * largeSeparatorsNum[i])
			}
		}
	}
	return
}
