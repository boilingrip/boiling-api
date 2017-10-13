package api

import (
	"errors"
	"fmt"
	"sort"
)

type LookupTable struct {
	forward map[string]int
	reverse map[int]string
	keys    []string
}

var ErrKeyNotFound = errors.New("key not found")

func BuildLookupTable(m map[int]string) LookupTable {
	t := LookupTable{
		forward: make(map[string]int),
		reverse: m,
	}

	for k, v := range m {
		if _, ok := t.forward[v]; ok {
			panic(fmt.Sprint("double value ", v))
		}

		t.forward[v] = k
		t.keys = append(t.keys, v)
	}

	sort.Strings(t.keys)

	return t
}

func (l LookupTable) Keys() []string {
	return l.keys
}

func (l LookupTable) LookUp(s string) (int, error) {
	v, ok := l.forward[s]
	if !ok {
		return -1, ErrKeyNotFound
	}

	return v, nil
}

func (l LookupTable) ReverseLookUp(i int) (string, error) {
	v, ok := l.reverse[i]
	if !ok {
		return "", ErrKeyNotFound
	}

	return v, nil
}

func (l LookupTable) Has(s string) bool {
	_, ok := l.forward[s]
	return ok
}

func (l LookupTable) HasReverse(i int) bool {
	_, ok := l.reverse[i]
	return ok
}
