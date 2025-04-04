package base62

import (
	"math"
	"strings"
)

var base62Str string

func MustInit(bs string) {
	if len(bs) == 0 {
		panic("need base string")
	}
	base62Str = bs
}

func Int2String(seq uint64) string {
	if seq == 0 {
		return string(base62Str[0])
	}
	bl := []byte{} // 存放结果

	for seq > 0 {
		mod := seq % 62
		div := seq / 62
		bl = append(bl, base62Str[mod])
		seq = div
	}

	return string(reverse(bl))
}

func reverse(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func string2Int(s string) (seq uint64) {
	bl := []byte(s)
	bl = reverse(bl)

	// 从右往左遍历
	for idx, b := range bl {
		base := math.Pow(62, float64(idx))
		seq += uint64(strings.Index(base62Str, string(b))) * uint64(base)
	}
	return seq
}
