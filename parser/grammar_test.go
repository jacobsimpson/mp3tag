package parser

import (
	"testing"

	"github.com/jacobsimpson/mp3tag/ast"
	"github.com/jacobsimpson/mp3tag/metadata"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input string
		want  ast.Expression
	}{
		{},
		{
			"title=abcfff",
			&ast.Equal{metadata.Title, "abcfff"},
		},
		{
			"    artist=thisistheartist    ",
			&ast.Equal{metadata.Artist, "thisistheartist"},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			assert := assert.New(t)

			got, err := Parse("", []byte(test.input))

			assert.NoError(err)
			assert.Equal(test.want, got)
		})
	}
}
