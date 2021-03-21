package cmd

import (
	"fmt"
	"os"

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
