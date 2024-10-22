package common

import (
	"math/rand"
	"strings"
	"time"
)

var obfuscationMap = map[int32][]rune{
	'A': {'A', 'a', '4', '@'}, 'B': {'B', 'b', '6', '8'}, 'C': {'C', 'c', '('},
	'D': {'D', 'd'}, 'E': {'E', 'e', '3'}, 'F': {'F', 'f', '1'},
	'G': {'G', 'g', '6', '9'}, 'H': {'H', 'h'}, 'I': {'I', 'i', '1', 'l', '!'},
	'J': {'J', 'j'}, 'K': {'K', 'k'}, 'L': {'L', 'l', '1', 'I', '!'},
	'M': {'M', 'm'}, 'N': {'N', 'n'}, 'O': {'O', 'o', '0', '#'},
	'P': {'P', 'p'}, 'Q': {'Q', 'q', '9'}, 'R': {'R', 'r'},
	'S': {'S', 's', '5', '$'}, 'T': {'T', 't', '7'}, 'U': {'U', 'u'},
	'V': {'V', 'v'}, 'W': {'W', 'w'}, 'X': {'X', 'x'},
	'Y': {'Y', 'y'}, 'Z': {'Z', 'z', '2', '?'}, '0': {'0', 'o', 'O', '#'},
	'1': {'1', 'l', 'I'}, '2': {'2', 'z', 'Z', '?'}, '3': {'3', 'e', 'E'},
	'4': {'4', 'a', 'A'}, '5': {'5', 's', 'S'}, '6': {'6', 'G', 'b'},
	'7': {'7', 'T'}, '8': {'8', 'B', '&'}, '9': {'9', 'g'},
}

var randGenerator *rand.Rand

func init() {
	src := rand.NewSource(time.Now().UnixNano())
	randGenerator = rand.New(src)
}

func ObfuscateText(input string) string {
	var result strings.Builder

	for _, char := range strings.ToUpper(input) {
		if options, exists := obfuscationMap[char]; exists {
			result.WriteRune(options[randGenerator.Intn(len(options))])
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}
