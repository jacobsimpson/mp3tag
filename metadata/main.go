package metadata

import (
	"fmt"

	id3 "github.com/mikkyang/id3-go"
)

type Name string

const (
	Album  Name = "album"
	Artist      = "artist"
	Title       = "title"
	Year        = "year"
	Genre       = "genre"
)

func AsName(name string) (Name, bool) {
	switch name {
	case "album":
		return Album, true
	case "artist":
		return Artist, true
	case "title":
		return Title, true
	case "year":
		return Year, true
	case "genre":
		return Genre, true
	}
	return "", false
}

type Tags struct {
	Album  string
	Artist string
	Title  string
	Year   string
	Genre  string
}

func (t *Tags) Value(name Name) string {
	switch name {
	case Album:
		return t.Album
	case Artist:
		return t.Artist
	case Title:
		return t.Title
	case Year:
		return t.Year
	case Genre:
		return t.Genre
	}
	// As long as this switch statement above is exhaustive, (which it is
	// supposed to be), this return statement will never be reached.
	return t.Genre
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
