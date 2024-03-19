package z85

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	testCases := []struct {
		name      string
		input     []byte
		expected  string
		expectErr bool
	}{
		{
			name:     "Valid",
			input:    []byte{0x86, 0x4F, 0xD2, 0x6F, 0xB5, 0x59, 0xF7, 0x5B},
			expected: "HelloWorld",
		},
		{
			name:      "Not divisible by 4",
			input:     []byte{0x01, 0x02, 0x03, 0x04, 0x05},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := Encode(tc.input)

			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		expected  []byte
		expectErr bool
	}{
		{
			name:     "Valid",
			input:    "HelloWorld",
			expected: []byte{0x86, 0x4F, 0xD2, 0x6F, 0xB5, 0x59, 0xF7, 0x5B},
		},
		{
			name:      "Not divisible by 5",
			input:     "four",
			expectErr: true,
		},
		{
			name:      "Invalid character",
			input:     "~~~~~",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := Decode(tc.input)

			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestEncodeDecode(t *testing.T) {
	testCases := []struct {
		input string
	}{
		{
			input: "JTKVSB%%)wK0E.X)V>+}o?pNmC{O&4W4b!Ni{Lh6",
		},
		{
			input: "ValidStringHere",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			decoded, err := Decode(tc.input)
			require.Nil(t, err)
			encoded, err := Encode(decoded)

			require.Nil(t, err)

			require.Equal(t, tc.input, encoded)
		})
	}
}
