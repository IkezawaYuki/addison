package main

import "unicode"

func IsPalindrome(s string) bool {
	var letter []rune
	for _, r := range s {
		if unicode.IsLetter(r) {
			letter = append(letter, r)
		}
	}
	for i := range letter {
		if letter[i] != letter[len(letter)-1-i] {
			return false
		}
	}
	return true
}
