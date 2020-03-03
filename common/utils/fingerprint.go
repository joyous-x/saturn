package utils

import (
	"fmt"
	"hash/crc32"
	"net/url"
	"strings"
)

const table uint32 = 0xEDB88321

func replaceCode(text string) string {
	var replaceCode = [12]string{"~", "!", "*", "%", "(", ")", "_", "-", ".", "'", " ", "&"}
	var replaceText = [12]string{"bab72b930fb3", "8476f727b721", "498db1d1d044", "46d30839d6bb", "bff45ea8c1d9", "6ca913c5a102", "245d5adc331c", "d7415efb3d4a", "5254000cc15a", "1d582bd01ea0", "dfefcf4564c9", "bab721030fb3"}

	for a := 0; a < 12; a++ {
		if strings.Contains(text, replaceCode[a]) {
			text = strings.Replace(text, replaceCode[a], replaceText[a], -1)
		}
	}
	return text
}

func Fingerprint(text string) string {
	crc32q := crc32.MakeTable(table)
	t := replaceCode(text)
	ts := url.QueryEscape(t)
	sigts := crc32.Checksum([]byte(ts), crc32q)
	return fmt.Sprintf("%x", sigts)
}
