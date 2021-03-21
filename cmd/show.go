package cmd

import (
	"fmt"
	"os"

	id3 "github.com/mikkyang/id3-go"
	"github.com/spf13/cobra"
)

type showFlagsStruct struct {
	album  bool
	artist bool
	genre  bool
	title  bool
	year   bool
}

var showFlags showFlagsStruct

func init() {
	showCmd.Flags().BoolVarP(&showFlags.album, "album", "b", false, "")
	showCmd.Flags().BoolVarP(&showFlags.artist, "artist", "a", false, "")
	showCmd.Flags().BoolVarP(&showFlags.genre, "genre", "g", false, "")
	showCmd.Flags().BoolVarP(&showFlags.title, "title", "t", false, "")
	showCmd.Flags().BoolVarP(&showFlags.year, "year", "y", false, "")
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
		showFile(filename)
	}
}

func showFile(filename string) {
	fmt.Printf("Working on %q\n", filename)
	mp3File, err := id3.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open %q: %+v\n", filename, err)
		os.Exit(1)
	}
	defer mp3File.Close()

	if showFlags.album {
		fmt.Println(mp3File.Album())
	}
	if showFlags.artist {
		fmt.Println(mp3File.Artist())
	}
	if showFlags.title {
		fmt.Println(mp3File.Title())
	}
	if showFlags.year {
		fmt.Println(mp3File.Year())
	}
	if showFlags.genre {
		fmt.Println(mp3File.Genre())
	}

}
