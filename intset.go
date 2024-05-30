package main

import (
	"bytes"
	"fmt"
	"sort"
)

type IntSet struct {
	words []uint64
}

func (s *IntSet) Has(x int) bool {
	word, bit := x/64, uint(x%64)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

func (s *IntSet) Add(x int) {
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", 64*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

func Example_one() {
	var x, y IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)
	fmt.Println(x.String())

	y.Add(9)
	y.Add(42)
	fmt.Println(y.String())
}

type MapIntSet struct {
	set map[int]bool
}

func (m *MapIntSet) Has(x int) bool {
	if m.set == nil {
		return false
	}
	return m.set[x]
}

func (m *MapIntSet) Add(x int) {
	if m.set == nil {
		m.set = make(map[int]bool)
	}
	m.set[x] = true
}

func (m *MapIntSet) UnionWith(t *MapIntSet) {
	if t.set == nil {
		return
	}
	if m.set == nil {
		m.set = make(map[int]bool)
	}
	for x, b := range t.set {
		if b {
			m.set[x] = true
		}
	}
}

func (s *MapIntSet) String() string {
	if s.set == nil {
		return "{ }"
	}
	ints := make([]int, 0, len(s.set))
	for x, v := range s.set {
		if v {
			ints = append(ints, x)
		}
	}
	sort.Ints(ints)

	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, x := range ints {
		if i != 0 {
			buf.WriteByte(' ')
		}
		_, _ = fmt.Fprintf(&buf, "%d", x)
	}
	buf.WriteByte('}')
	return buf.String()
}
