package textutil

import (
	"encoding/hex"
	"strings"
)

const runeOrder = " aäąàbcćdeęfghijklłmnńoöópqrsßśtuüvwxyzźż" +
	"абвгґдеёжзіийклмнопрстуўфхцчшщъыьэюя"

var runeOrderMap = map[rune]byte{}

func init() {
	for i, r := range runeOrder {
		runeOrderMap[r] = byte(i)
	}
}

func getRuneOrder(r rune) byte {
	if b, ok := runeOrderMap[r]; ok {
		return b
	} else {
		return 255
	}
}

func CreateSortKey(s string) string {
	s = strings.ToLower(s)
	sr := []rune(s)
	buff := make([]byte, 0, len(sr)+1)
	for _, r := range sr {
		if r == '’' || r == '.' {
			continue
		}
		buff = append(buff, getRuneOrder(r))
	}

	return hex.EncodeToString(buff)
}

func RuneOrderCmp(a, b string) int {
	ar, br := []rune(a), []rune(b)

	for {
		if len(ar) == 0 && len(br) == 0 {
			return 0
		}
		if len(ar) == 0 {
			return -1
		}
		if len(br) == 0 {
			return 1
		}
		result := int(getRuneOrder(ar[0])) - int(getRuneOrder(br[0]))
		if result == 0 {
			ar, br = ar[1:], br[1:]
			continue
		}
		return result
	}
}
