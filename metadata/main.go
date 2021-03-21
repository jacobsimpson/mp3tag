package metadata

import (
	"fmt"

	id3 "github.com/mikkyang/id3-go"
)

type Tags struct {
	Album  string
	Artist string
	Title  string
	Year   string
	Genre  string
}

func (t *Tags) Value(name string) (string, bool) {
	switch name {
	case "album":
		return t.Album, true
	case "artist":
		return t.Artist, true
	case "title":
		return t.Title, true
	case "year":
		return t.Year, true
	case "genre":
		return t.Genre, true
	}
	return "", false
}

func ReadTags(filename string) (*Tags, error) {
	mp3File, err := id3.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open file %q: %+v", filename, err)
	}
	defer mp3File.Close()

	return &Tags{
		Album:  mp3File.Album(),
		Artist: mp3File.Artist(),
		Title:  mp3File.Title(),
		Year:   mp3File.Year(),
		Genre:  mp3File.Genre(),
	}, nil
}
