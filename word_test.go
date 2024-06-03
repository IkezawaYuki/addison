package main

import (
	"math/rand"
	"testing"
	"unicode"
)

func TestIsPalindrome(t *testing.T) {
	var tests = []struct {
		input string
		want  bool
	}{
		{"", true},
		{"a", true},
		{"aa", true},
		{"ab", false},
		{"kayak", true},
		{"detartrated", true},
		{"A man, a plan, a canal: Panama", true},
		{"Evil I did dwell; lewd did I live.", true},
		{"Able was I ere I saw Elba", true},
		{"été", true},
		{"Et se resservir, ivresse reste.", true},
		{"palindrome", false}, // non-palindrome
		{"desserts", false},   // semi-palindrome
	}
	for _, test := range tests {
		if got := IsPalindrome(test.input); got != test.want {
			t.Errorf("IsPalindrome(%q) = %v", test.input, got)
		}
	}
}

func BenchmarkIsPalindrome(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsPalindrome("A man, a plan, a canal: Panama")
	}
}

func randomPalindrome(rng *rand.Rand) string {
	n := rng.Intn(25) + 2
	runes := make([]rune, n)
	for i := 0; i < n; i++ {
		r := rune(rng.Intn(0x1000))
		if unicode.IsLetter(r) {
			runes[i] = r
			i++
		}
	}
	for i := 0; i < (n+1)/2; i++ {
		if runes[i] == runes[n-1-i] {
			return string(runes)
		}
	}
	return ""
}
