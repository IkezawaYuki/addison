package main

import (
	"bytes"
	"testing"
)

func TestCharcount(t *testing.T) {
	for _, test := range []struct {
		bytes   []byte
		counts  map[rune]int
		utflen  []int
		invalid int
	}{
		{
			[]byte("こんにちは世界"),
			map[rune]int{'こ': 1, 'ん': 1, 'に': 1, 'ち': 1, 'は': 1, '、': 1, '世': 1, '界': 1},
			[]int{0, 0, 0, 8, 0},
			0,
		},
		{
			[]byte("Hello, World"),
			map[rune]int{'H': 1, 'e': 1, 'l': 3, 'o': 2, ',': 1, ' ': 1, 'W': 1, 'r': 1, 'd': 1},
			[]int{0, 12, 0, 0, 0},
			0,
		},
		{
			[]byte("Hello, World\300"), // the last byte is invalid
			map[rune]int{'H': 1, 'e': 1, 'l': 3, 'o': 2, ',': 1, ' ': 1, 'W': 1, 'r': 1, 'd': 1},
			[]int{0, 12, 0, 0, 0},
			1,
		},
	} {
		counts, utflen, invalid, err := charcount(bytes.NewBuffer(test.bytes))
		if err != nil {
			t.Error(err)
			continue
		}
		if len(counts) != len(test.counts) {
			t.Errorf("len(counts): got %d, want %d", len(counts), len(test.counts))
			continue
		}
		for k, v := range test.counts {
			count, ok := counts[k]
			if !ok {
				t.Errorf("%c is not included,", k)
				continue
			}
			if count != v {
				t.Errorf("count for %c is %d, but want %d\n", k, count, v)
			}
		}
		if len(utflen) != len(test.utflen) {
			t.Errorf("len(utflen): got %d, want %d", len(utflen), len(test.utflen))
			continue
		}
		for i := 0; i < len(utflen); i++ {
			if utflen[i] != test.utflen[i] {
				t.Errorf("utflen[%d]: got %d, want %d", i, utflen[i], test.utflen[i])
				continue
			}
		}
		if invalid != test.invalid {
			t.Errorf("invalid: got %d, want %d", invalid, test.invalid)
			continue
		}
	}
}
