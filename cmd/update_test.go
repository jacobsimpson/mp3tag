package cmd

import (
	"testing"

	"github.com/jacobsimpson/mp3tag/metadata"
	"github.com/stretchr/testify/assert"
)

func TestExpandRenameFormat(t *testing.T) {
	tests := []struct {
		format string
		tags   *metadata.Tags
		want   string
	}{
		{
			"{title}.mp3",
			&metadata.Tags{Title: "this-is-the-title"},
			"this-is-the-title.mp3",
		},
		{
			"293ab{title}-file.mp3",
			&metadata.Tags{Title: "this-is-the-title"},
			"293abthis-is-the-title-file.mp3",
		},
		{
			"{title}this{title}xy{title}.mp3",
			&metadata.Tags{Title: "mytitle"},
			"mytitlethismytitlexymytitle.mp3",
		},
		{
			"{album}",
			&metadata.Tags{Title: "mytitle"},
			"",
		},
		{
			"thisisafilename",
			&metadata.Tags{Title: "mytitle"},
			"thisisafilename",
		},
		{
			"1{artist}2",
			&metadata.Tags{Artist: "mytitle"},
			"1mytitle2",
		},
	}

	for _, test := range tests {
		t.Run(test.format, func(t *testing.T) {
			assert := assert.New(t)
			got, err := expandRenameFormat(test.format, test.tags)
			assert.NoError(err)
			assert.Equal(test.want, got)
		})
	}
}
