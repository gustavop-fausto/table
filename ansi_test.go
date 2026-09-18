package table

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CSIStrip(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "hello world!",
			want:  "hello world!",
		},
		{
			input: "\x1b[37mhello \x1b[31mworld!",
			want:  "hello world!",
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, test.want, newANSI(test.input).Strip())
		})
	}
}

func Test_ANSILen(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{
			input: "hello world!",
			want:  12,
		},
		{
			input: "\x1b[37mhello \x1b[31mworld!",
			want:  12,
		},
		{
			input: "🔥 unicode 🔥 characters 🔥",
			want:  27,
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, test.want, newANSI(test.input).Len())
		})
	}
}

func TestNewANSI_Multibyte(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  ansiBlob
	}{
		{
			name:  "multibyte character before ANSI sequence",
			input: "Olá, \x1b[31mMundo",
			want: ansiBlob{
				{value: "Olá, ", style: ""},
				{value: "Mundo", style: "\x1b[31m"},
			},
		},
		{
			name:  "multibyte characters between ANSI sequence",
			input: "\x1b[31mCafé \x1b[32m3 Corações",
			want: ansiBlob{
				{value: "Café ", style: "\x1b[31m"},
				{value: "3 Corações", style: "\x1b[32m"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := newANSI(test.input)

			if !slices.Equal(test.want, got) {
				t.Errorf("expected %v, got %v", test.want, got)
			}
		})
	}
}
