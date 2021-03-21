package cmd

import (
	"fmt"
	"os"

	id3 "github.com/mikkyang/id3-go"
	"github.com/spf13/cobra"
)

var album bool
var artist bool
var genre bool
var title bool
var year bool

func init() {
	showCmd.Flags().BoolVarP(&album, "album", "b", false, "")
	showCmd.Flags().BoolVarP(&artist, "artist", "a", false, "")
	showCmd.Flags().BoolVarP(&genre, "genre", "g", false, "")
	showCmd.Flags().BoolVarP(&title, "title", "t", false, "")
	showCmd.Flags().BoolVarP(&year, "year", "y", false, "")
	rootCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the requested attributes of the specified files.",
	Long:  "Show the requested attributes of the specified files.",
	Args:  cobra.MinimumNArgs(1),
	Run:   showCmdRun,
}

func showCmdRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("No files specified.")
		return
	}

	for _, filename := range args {

		fmt.Printf("Working on %q\n", filename)
		mp3File, err := id3.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to open %q: %+v\n", filename, err)
			os.Exit(1)
		}
		defer mp3File.Close()

		if album {
			fmt.Println(mp3File.Album())
		}
		if artist {
			fmt.Println(mp3File.Artist())
		}
		if title {
			fmt.Println(mp3File.Title())
		}
		if year {
			fmt.Println(mp3File.Year())
		}
		if genre {
			fmt.Println(mp3File.Genre())
		}

	}
}
