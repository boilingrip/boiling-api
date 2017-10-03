package api

import (
	"errors"
	"sort"
	"strings"
	"unicode"
)

// prepareTags prepares tags.
// Tags are brought to lower case, sorted and deduplicated.
// If any tag contains whitespace, an error is returned.
func prepareTags(a []string) ([]string, error) {
	if len(a) == 0 {
		return a, nil
	}

	sort.Strings(a)
	var last string
	out := a[:0]
	for i, s := range a {
		tmp := strings.Map(func(r rune) rune {
			if unicode.IsSpace(r) {
				return -1
			}
			return r
		}, s)

		if tmp != s {
			//hacky, but works for now
			return nil, errors.New("tag contained whitespace")
		}

		s = strings.ToLower(sanitizeString(s))

		if i == 0 {
			last = s
			out = append(out, s)
			continue
		}
		if s != last {
			out = append(out, s)
			last = s
		}
	}

	return out, nil
}
