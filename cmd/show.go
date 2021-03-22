package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

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
	showCmd.Flags().BoolVarP(&showFlags.album, "album", "b", false, "show the album tag")
	showCmd.Flags().BoolVarP(&showFlags.artist, "artist", "a", false, "show the artist tag")
	showCmd.Flags().BoolVarP(&showFlags.genre, "genre", "g", false, "show the genre tag")
	showCmd.Flags().BoolVarP(&showFlags.title, "title", "t", false, "show the title tag")
	showCmd.Flags().BoolVarP(&showFlags.year, "year", "y", false, "show the year tag")
	rootCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:   "show [flags] file+",
	Short: "Show the requested tags of the specified files.",
	Long: strings.TrimSpace(`
Show the requested tags of the specified files. If no tags are specified, all
tags will be displayed.`),
	Example: `  Show all the tags for all the .mp3 files:
  	show *.mp3

  Show the artist and genre tags for podcast.mp3:
  	show --artist --genre podcast.mp3`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one file")
		}
		return nil
	},
	Run: showCmdRun,
}

func showCmdRun(cmd *cobra.Command, args []string) {
	// If no attributes were specified, then all attributes will be shown.
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

	fmt.Printf("======== %s ======== \n", filename)
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
