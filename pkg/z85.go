// Package z85 provides basic encode and decode functions for Z85 as described here:
// https://rfc.zeromq.org/spec/32/
package z85

import (
	"fmt"
	"strings"
)

// Z85LookupTable is the lookup table for the base85 encoding.
// The encoding and decoding SHALL use this representation for each base-85 value from zero to 84
var Z85LookupTable = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.-:+=^!/*?&<>()[]{}@%$#"

// Encode takes a binary frame and encodes it as printable ASCII string.
func Encode(data []byte) (string, error) {
	// The binary frame SHALL have a length that is divisible by 4 with no remainder
	if len(data)%4 != 0 {
		return "", fmt.Errorf("binary frame is not divisible by 4")
	}

	/*
		To encode a frame, an implementation SHALL take four octets at a time from the binary frame and convert them into five printable characters

		The four octets SHALL be treated as an unsigned 32-bit integer in network byte order (big endian)
	*/
	var result strings.Builder
	for i := 0; i < len(data); i += 4 {
		// Take the 4 octets from the binary frame and convert them to a uint32
		fourOctetsInt := uint32(data[i])<<24 | uint32(data[i+1])<<16 | uint32(data[i+2])<<8 | uint32(data[i+3])

		var characters [5]rune
		// The five characters SHALL be output from most significant to least significant (big endian)
		for j := 4; j >= 0; j-- {
			// Continually divide by 85 to get base 85 of something.
			remainder := fourOctetsInt % 85
			characters[j] = rune(Z85LookupTable[remainder])
			fourOctetsInt /= 85
		}

		for _, c := range characters {
			result.WriteRune(c)
		}
	}
	return result.String(), nil
}

// Decode takes a printable ASCII string and decodes it into binary frame.
func Decode(data string) ([]byte, error) {
	// The string frame SHALL have a length that is divisible by 5 with no remainder
	if len(data)%5 != 0 {
		return nil, fmt.Errorf("string frame is not divisible by 5")
	}

	/*
		To decode a string, an implementation SHALL take five characters at a time from the string and convert them into four octets of data representing a 32-bit unsigned integer in network byte order (big endian)

		The five characters SHALL each be converted into a value 0 to 84, and accumulated by multiplication by 85, from most to least significant.
	*/
	result := make([]byte, len(data)/5*4)
	for i := 0; i < len(data); i += 5 {
		fiveChars := data[i : i+5]

		var fiveCharsInt uint32
		for j, char := range fiveChars {
			charIndex := strings.Index(Z85LookupTable, string(char))
			if charIndex == -1 {
				return nil, fmt.Errorf("invalid character in string frame: %c", char)
			}
			// Multiply by 85 for each iteration (except the last one)
			if j < 4 {
				// For the first four characters
				fiveCharsInt = (fiveCharsInt + uint32(charIndex)) * 85
			} else {
				// For the last character, just add it
				fiveCharsInt += uint32(charIndex)
			}
		}
		// Calculates byte values from the 32-bit integer
		// Compute the index in the result slice directly based on the loop counter
		index := (i / 5) * 4
		result[index] = byte(fiveCharsInt >> 24 & 0xFF)
		result[index+1] = byte(fiveCharsInt >> 16 & 0xFF)
		result[index+2] = byte(fiveCharsInt >> 8 & 0xFF)
		result[index+3] = byte(fiveCharsInt & 0xFF)
	}

	return result, nil
}
