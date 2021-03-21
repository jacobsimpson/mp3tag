package cmd

import (
	"fmt"
	"os"

	id3 "github.com/mikkyang/id3-go"
	"github.com/spf13/cobra"
)

type updateFlagsStruct struct {
	album        string
	artist       string
	genre        string
	title        string
	year         string
	renameFormat string
}

var updateFlags updateFlagsStruct

func init() {
	updateCmd.Flags().StringVarP(&updateFlags.album, "album", "b", "", "")
	updateCmd.Flags().StringVarP(&updateFlags.artist, "artist", "a", "", "")
	updateCmd.Flags().StringVarP(&updateFlags.genre, "genre", "g", "", "")
	updateCmd.Flags().StringVarP(&updateFlags.title, "title", "t", "", "")
	updateCmd.Flags().StringVarP(&updateFlags.year, "year", "y", "", "")
	updateCmd.Flags().StringVarP(&updateFlags.renameFormat, "rename", "r", "", "")
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the requested attributes of the specified files.",
	Long:  "Update the requested attributes of the specified files.",
	Args:  cobra.MinimumNArgs(1),
	Run:   updateCmdRun,
}

func updateCmdRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("No files specified.")
		return
	}

	for _, filename := range args {
		if tagUpdates() {
			updateFile(filename)
		}
		if updateFlags.renameFormat != "" {
			renameFile(filename)
		}
	}

}

func tagUpdates() bool {
	return updateFlags.album != "" ||
		updateFlags.artist != "" ||
		updateFlags.title != "" ||
		updateFlags.year != "" ||
		updateFlags.genre != ""
}

func updateFile(filename string) {

	fmt.Printf("Working on %q\n", filename)
	mp3File, err := id3.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open %q: %+v\n", filename, err)
		return
	}
	defer mp3File.Close()

	if updateFlags.album != "" {
		mp3File.SetAlbum(updateFlags.artist)
	}
	if updateFlags.artist != "" {
		mp3File.SetArtist(updateFlags.artist)
	}
	if updateFlags.title != "" {
		mp3File.SetTitle(updateFlags.title)
	}
	if updateFlags.year != "" {
		mp3File.SetYear(updateFlags.year)
	}
	if updateFlags.genre != "" {
		mp3File.SetGenre(updateFlags.genre)
	}

}

func renameFile(src string) error {
	tags, err := readTags(src)
	if err != nil {
		return err
	}
	dst, err := expandRenameFormat(updateFlags.renameFormat, tags)
	if err != nil {
		return err
	}
	if err := os.Rename(src, dst); err != nil {
		return fmt.Errorf("Unable to rename %q to %q: %+v", src, dst, err)
	}
	return nil
}

type parseState int

const (
	initial parseState = iota
	variable
)

func expandRenameFormat(renameFormat string, tags *Tags) (string, error) {
	result := ""
	name := ""
	state := initial

	for _, r := range renameFormat {
		switch state {
		case initial:
			switch r {
			case '{':
				name = ""
				state = variable
			default:
				result += string(r)
			}
		case variable:
			switch r {
			case '}':
				v, err := tags.Value(name)
				if err != nil {
					return "", err
				}
				result += v
				state = initial
			default:
				name += string(r)
			}
		}
	}
	return result, nil
}

type Tags struct {
	Album  string
	Artist string
	Title  string
	Year   string
	Genre  string
}

func (t *Tags) Value(name string) (string, error) {
	switch name {
	case "album":
		return t.Album, nil
	case "artist":
		return t.Artist, nil
	case "title":
		return t.Title, nil
	case "year":
		return t.Year, nil
	case "genre":
		return t.Genre, nil
	}
	return "", fmt.Errorf("no such attribute %q", name)
}

func readTags(filename string) (*Tags, error) {
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
