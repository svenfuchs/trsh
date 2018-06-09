package util

import (
	// "bytes"
	"fmt"
	"io"
	"strings"
)

func HasPrefix(str string, prefix ...string) bool {
	for _, p := range prefix {
		if strings.HasPrefix(str, p) {
			return true
		}
	}
	return false
}

func Split(str string, delim string) []string {
	strs := []string{}
	for _, s := range strings.Split(str, delim) {
		strs = append(strs, Strip(s))
	}
	return strs
}

func Strip(str string) string {
	return StripLeft(StripRight(str))
}

func StripLeft(str string) string {
	return strings.TrimPrefix(str, " ")
}

func StripRight(str string) string {
	return strings.TrimSuffix(str, " ")
}

func StripPrefix(str string, prefix ...string) string {
	str = StripLeft(str)
	for _, s := range prefix {
		str = strings.TrimPrefix(str, s)
		str = StripLeft(str)
	}
	return str
}

func StripSuffix(str string, suffix ...string) string {
	str = StripRight(str)
	for _, s := range suffix {
		str = strings.TrimSuffix(str, s)
		str = StripRight(str)
	}
	return str
}

func HasSuffix(str string, strs ...string) bool {
	for _, s := range strs {
		if strings.HasSuffix(str, s) {
			return true
		}
	}
	return false
}

func LastPart(str string, delim string) string {
	parts := strings.Split(str, " ")
	return parts[len(parts)-1]
}

func FindMatch(strs []string, str string) string {
	if str == "" {
		return ""
	}
	for _, s := range strs {
		if strings.HasPrefix(s, str) {
			return s
		}
	}
	return ""
}

func MatchStrs(strs []string, str string) []string {
	res := []string{}
	if str == "" {
		return res
	}
	for _, s := range strs {
		if strings.HasPrefix(s, str) {
			res = append(res, s)
		}
	}
	return res
}

func MatchesStrs(strs []string, str string) bool {
	if str == "" {
		return false
	}
	for _, s := range strs {
		if strings.HasPrefix(s, str) {
			return true
		}
	}
	return false
}

func Unique(strs []string) []string {
	u := make([]string, 0, len(strs))
	m := make(map[string]bool)

	for _, s := range strs {
		if _, ok := m[s]; !ok {
			m[s] = true
			u = append(u, s)
		}
	}

	return u
}

func Includes(strs []string, str string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}

func Except(strs []string, other ...string) []string {
	res := []string{}
	for _, s := range strs {
		if !Includes(other, s) {
			res = append(res, s)
		}
	}
	return res
}

func Combinations(strs []string) []string {
	c := []string{}
	for _, s := range strs {
		c = append(c, s)
		for _, r := range Combinations(Except(strs, s)) {
			c = append(c, s+" "+r)
		}
	}
	return c
}

func FirstOr(strs []string, d string) string {
	if len(strs) > 0 {
		return strs[0]
	}
	return d
}

func Times(n int, f func(int)) {
	for i := 0; i < n; i++ {
		f(i)
	}
}

func Read(r io.Reader) []byte {
	// TODO why does this panic with "slice bounds out of range"??
	// buf := new(bytes.Buffer)
	// buf.ReadFrom(r)
	// return buf.Bytes()
	b := make([]byte, 100000)
	c, _ := r.Read(b)
	return b[:c]
}

func Write(w io.Writer, s string) {
	fmt.Fprintln(w, s)
}
