package service

import "testing"

func TestSlugify(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Hello World", "hello-world"},
		{"My Great Post!", "my-great-post"},
		{"  leading and trailing  ", "leading-and-trailing"},
		{"UPPER CASE", "upper-case"},
		{"special---chars!!!here", "special-chars-here"},
		{"hello---world", "hello-world"},
		{"café au lait", "caf-au-lait"},
		{"123 Numbers", "123-numbers"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := slugify(tt.input)
			if got != tt.want {
				t.Errorf("slugify(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
