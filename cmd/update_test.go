package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandRenameFormat(t *testing.T) {
	tests := []struct {
		format string
		tags   *Tags
		want   string
	}{
		{
			"{title}.mp3",
			&Tags{Title: "this-is-the-title"},
			"this-is-the-title.mp3",
		},
		{
			"293ab{title}-file.mp3",
			&Tags{Title: "this-is-the-title"},
			"293abthis-is-the-title-file.mp3",
		},
		{
			"{title}this{title}xy{title}.mp3",
			&Tags{Title: "mytitle"},
			"mytitlethismytitlexymytitle.mp3",
		},
		{
			"{album}",
			&Tags{Title: "mytitle"},
			"",
		},
		{
			"thisisafilename",
			&Tags{Title: "mytitle"},
			"thisisafilename",
		},
		{
			"1{artist}2",
			&Tags{Artist: "mytitle"},
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
