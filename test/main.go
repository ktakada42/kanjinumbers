package test

import (
	"fmt"
	"os"
	"strconv"

	"kanjinumbers/pkg/server/handler"
)

func main() {
	for i := handler.Min; i < handler.Max; i++ {
		n := strconv.Itoa(i)
		kanji := handler.ConvertNumberToKanji(n)
		num, _ := handler.ConvertKanjiToNumber(kanji)
		if i == num {
		} else {
			fmt.Printf("%d: NG :(\n", i)
			os.Exit(1)
		}
	}
}
