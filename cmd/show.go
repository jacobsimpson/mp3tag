package cmd

import (
	"fmt"
	"os"

	id3 "github.com/mikkyang/id3-go"
	"github.com/spf13/cobra"
)

type showFlagsStruct struct {
	all    bool
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

	showFlags.all = !(showFlags.album || showFlags.artist || showFlags.genre || showFlags.title || showFlags.year)
	for _, filename := range args {
		showFile(filename)
	}
}

func showFile(filename string) {
	mp3File, err := id3.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open %q: %+v\n", filename, err)
		os.Exit(1)
	}
	defer mp3File.Close()

	fmt.Println("========")
	if showFlags.album || showFlags.all {
		fmt.Printf("Album:  %s\n", mp3File.Album())
	}
	if showFlags.artist || showFlags.all {
		fmt.Printf("Artist: %s\n", mp3File.Artist())
	}
	if showFlags.title || showFlags.all {
		fmt.Printf("Title:  %s\n", mp3File.Title())
	}
	if showFlags.year || showFlags.all {
		fmt.Printf("Year:   %s\n", mp3File.Year())
	}
	if showFlags.genre || showFlags.all {
		fmt.Printf("Genre:  %s\n", mp3File.Genre())
	}

}
