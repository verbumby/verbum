package textutil

import (
	"bufio"
	"bytes"
)

func GetDelimSplitFunc(delim string) bufio.SplitFunc {
	delimBytes := []byte(delim)

	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if i := bytes.Index(data, delimBytes); i >= 0 {
			return i + len(delimBytes), data[0:i], nil
		}

		if atEOF {
			return len(data), data, nil
		}

		return 0, nil, nil
	}
}
