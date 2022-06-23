package test

import (
	"strconv"
	"testing"

	"kanjinumbers/pkg/server/handler"
)

func TestMinToMax(t *testing.T) {
	for i := handler.Min; i < handler.Max; i++ {
		n := strconv.Itoa(i)
		kanji := handler.ConvertNumberToKanji(n)
		num, _ := handler.ConvertKanjiToNumber(kanji)
		if i == num {
		} else {
			t.Errorf("%d: NG :(", i)
		}
	}
}
