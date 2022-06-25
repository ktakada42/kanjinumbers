package test

import (
	"strconv"
	"testing"

	"kanjinumbers/pkg/server/handler"
)

func TestConvertValueCompare(t *testing.T) {
	t.Run("0", func(t *testing.T) {
		t.Parallel()
		i := 0
		n := strconv.Itoa(i)
		kanji := handler.ConvertNumberToKanji(n)
		num, _ := handler.ConvertKanjiToNumber(kanji)
		if i == num {
			t.Logf("\x1b[32m\nparameter %d: kanji: %s\n\x1b[0m", i, kanji)
		} else {
			t.Fatalf("\x1b[31m\nparameter %d: kanji: %s\n\x1b[0m", i, kanji)
		}
	})

	for i := 1; i <= 9999; i++ {
		i := i

		t.Run("1 to 9,999", func(t *testing.T) {
			t.Parallel()
			n := strconv.Itoa(i)
			kanji := handler.ConvertNumberToKanji(n)
			num, _ := handler.ConvertKanjiToNumber(kanji)
			if i == num {
				t.Logf("\x1b[32m\nparameter %d: kanji: %s\n\x1b[0m", i, kanji)
			} else {
				t.Fatalf("\x1b[31m\nparameter %d: kanji: %s\n\x1b[0m", i, kanji)
			}
		})

		t.Run("10,000 to 99,990,000", func(t *testing.T) {
			t.Parallel()
			i := i * 10000
			n := strconv.Itoa(i)
			kanji := handler.ConvertNumberToKanji(n)
			num, _ := handler.ConvertKanjiToNumber(kanji)
			if i == num {
				t.Logf("\x1b[32m\nparameter %d: kanji: %s\n\x1b[0m", i, kanji)
			} else {
				t.Fatalf("\x1b[31m\nparameter %d: kanji: %s\n\x1b[0m", i, kanji)
			}
		})

		t.Run("100,000,000 to 999,900,000", func(t *testing.T) {
			t.Parallel()
			i := i * 100000000
			n := strconv.Itoa(i)
			kanji := handler.ConvertNumberToKanji(n)
			num, _ := handler.ConvertKanjiToNumber(kanji)
			if i == num {
				t.Logf("\x1b[32m\nparameter %d: kanji: %s\n\x1b[0m", i, kanji)
			} else {
				t.Fatalf("\x1b[31m\nparameter %d: kanji: %s\n\x1b[0m", i, kanji)
			}
		})

		t.Run("1,000,000,000,000 to 9,999,000,000,000,000", func(t *testing.T) {
			t.Parallel()
			i := i * 1000000000000
			n := strconv.Itoa(i)
			kanji := handler.ConvertNumberToKanji(n)
			num, _ := handler.ConvertKanjiToNumber(kanji)
			if i == num {
				t.Logf("\x1b[32m\nparameter %d: kanji: %s\n\x1b[0m", i, kanji)
			} else {
				t.Fatalf("\x1b[31m\nparameter %d: kanji: %s\n\x1b[0m", i, kanji)
			}
		})
	}
}
