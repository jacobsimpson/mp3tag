package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jacobsimpson/mp3tag/metadata"
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
	updateCmd.Flags().StringVarP(&updateFlags.album, "album", "b", "", "specify a new value for the album tag")
	updateCmd.Flags().StringVarP(&updateFlags.artist, "artist", "a", "", "specify a new value for the artist tag")
	updateCmd.Flags().StringVarP(&updateFlags.genre, "genre", "g", "", "specify a new value for the genre tag")
	updateCmd.Flags().StringVarP(&updateFlags.title, "title", "t", "", "specify a new value for the title tag")
	updateCmd.Flags().StringVarP(&updateFlags.year, "year", "y", "", "specify a new value for the year tag")
	updateCmd.Flags().StringVarP(&updateFlags.renameFormat, "rename", "r", "", "rename the file using tag values")
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update [flags] file+",
	Short: "Update the requested attributes of the specified files.",
	Long: strings.TrimSpace(`
Update the requested attributes of the specified files. The values specified
will be set on all files specified. This behavior is a natural fit for
attributes like album, artist, or genre where many files may share the same
values. It is probably less interesting to set the same title on the entire
list of files.`),
	Example: `  Set 'John Doe' as the artist for all mp3 files:
  	update --artist="John Doe" *.mp3

  Set the artist to 'John Doe', the album to 'The Best Podcast', and the genre to
  'Podcast' for all mp3 files in the podcast directory:
  	update --artist="John Doe" --album="The Best Podcast" --genre=Podcast podcast/*.mp3

  Rename all .mp3 files in the podcast directory to '<artist> - <album> -
  <title>.mp3'. The actual name of the file will depend on the values each of
  those tags has for each file:
  	update --rename="{artist} - {album} - {title}.mp3" podcast/*.mp3`,
	Args: func(cmd *cobra.Command, args []string) error {
		any := updateFlags.album != "" ||
			updateFlags.artist != "" ||
			updateFlags.title != "" ||
			updateFlags.year != "" ||
			updateFlags.genre != "" ||
			updateFlags.renameFormat != ""
		if !any {
			return errors.New("requires at least one flag to update")
		}

		if len(args) < 1 {
			return errors.New("requires at least one file")
		}

		return nil
	},
	Run: updateCmdRun,
}

func updateCmdRun(cmd *cobra.Command, args []string) {
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

	fmt.Printf("Updating %q\n", filename)
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
	tags, err := metadata.ReadTags(src)
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

func expandRenameFormat(renameFormat string, tags *metadata.Tags) (string, error) {
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
				n, ok := metadata.AsName(name)
				if !ok {
					return "", fmt.Errorf("%q is not a valid tag name", name)
				}
				result += tags.Value(n)
				state = initial
			default:
				name += string(r)
			}
		}
	}
	return result, nil
}
