package dictimport

import (
	"encoding/hex"
	"strings"
)

const charOrder = " aäąàbcćdeęfghijklłmnńoöópqrsßśtuüvwxyzźż" +
	"абвгдеёжзіийклмнопрстуўфхцчшщъыьэюя"

var charOrderMap = map[rune]byte{}

func init() {
	for i, r := range charOrder {
		charOrderMap[r] = byte(i)
	}
}

func createSortKey(s string) string {
	s = strings.ToLower(s)
	sr := []rune(s)
	buff := make([]byte, 0, len(sr)+1)
	for _, r := range sr {
		if r == '’' || r == '.' {
			continue
		}
		if b, ok := charOrderMap[r]; ok {
			buff = append(buff, b)
		} else {
			buff = append(buff, 255)
		}
	}

	return hex.EncodeToString(buff)
}
