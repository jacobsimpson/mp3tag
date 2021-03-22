package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/jacobsimpson/mp3tag/ast"
	"github.com/jacobsimpson/mp3tag/metadata"
	"github.com/jacobsimpson/mp3tag/parser"
	"github.com/spf13/cobra"
)

type findFlagsStruct struct {
}

var findFlags findFlagsStruct

func init() {
	//updateCmd.Flags().StringVarP(&updateFlags.album, "album", "b", "", "")
	//updateCmd.Flags().StringVarP(&updateFlags.artist, "artist", "a", "", "")
	//updateCmd.Flags().StringVarP(&updateFlags.genre, "genre", "g", "", "")
	//updateCmd.Flags().StringVarP(&updateFlags.title, "title", "t", "", "")
	//updateCmd.Flags().StringVarP(&updateFlags.year, "year", "y", "", "")
	//updateCmd.Flags().StringVarP(&updateFlags.renameFormat, "rename", "r", "", "")
	rootCmd.AddCommand(findCmd)
}

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Find files that match the query.",
	Long: `Find files that match the query.

	Queries take the form of:

	title="abc"

	title=""
	`,
	Args: cobra.MinimumNArgs(2),
	Run:  findCmdRun,
}

func findCmdRun(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Println("No files specified.")
		return
	}

	query, err := parser.Parse("", []byte(args[0]), parser.Debug(false))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse query: %+v\n", err)
		return
	}

	if err := execute(query.(ast.Expression), args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to execute query: %+v\n", err)
		return
	}
}

func execute(query ast.Expression, files []string) error {
	for _, file := range files {
		info, err := os.Stat(file)
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s does not exist.\n", file)
			continue
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "%s could not be opened: %+v\n", file, err)
			continue
		}
		if info.IsDir() {
			fmt.Fprintf(os.Stderr, "%s is a directory.\n", file)
			continue
		}

		if is, err := isMP3(file); err != nil {
			fmt.Fprintf(os.Stderr, "%s could not be opened: %+v\n", file, err)
			continue
		} else if !is {
			fmt.Fprintf(os.Stderr, "%s is not an MP3 file.\n", file)
			continue
		}

		tags, err := metadata.ReadTags(file)
		if err != nil {
			return err
		}

		if m, err := match(query, tags); m {
			fmt.Println(file)
		} else if err != nil {
			return err
		}
	}
	return nil
}

func match(query ast.Expression, tags *metadata.Tags) (bool, error) {
	switch q := query.(type) {
	case *ast.Equal:
		return equal(q, tags)
	case *ast.Has:
		return has(q, tags)
	}
	return false, nil
}

func equal(e *ast.Equal, tags *metadata.Tags) (bool, error) {
	return strings.ToLower(tags.Value(e.LHS)) == strings.ToLower(e.RHS), nil
}

func has(e *ast.Has, tags *metadata.Tags) (bool, error) {
	return strings.Contains(strings.ToLower(tags.Value(e.LHS)), strings.ToLower(e.RHS)), nil
}

// Attempt to recognize an MP3 file using magic numbers.
// https://en.wikipedia.org/wiki/List_of_file_signatures
var mp3MagicNumber1 = []byte{0xFF, 0xFB}
var mp3MagicNumber2 = []byte{0xFF, 0xF3}
var mp3MagicNumber3 = []byte{0xFF, 0xF2}
var id3MagicNumber = []byte{0x49, 0x44, 0x33}

func isMP3(filename string) (bool, error) {
	f, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer f.Close()

	buffer := make([]byte, 3)
	if _, err := f.Read(buffer); err != nil {
		return false, err
	}
	if bytes.Compare(mp3MagicNumber1, buffer[0:2]) == 0 {
		return true, nil
	} else if bytes.Compare(mp3MagicNumber2, buffer[0:2]) == 0 {
		return true, nil
	} else if bytes.Compare(mp3MagicNumber3, buffer[0:2]) == 0 {
		return true, nil
	} else if bytes.Compare(id3MagicNumber, buffer) == 0 {
		return true, nil
	}
	return false, nil
}
