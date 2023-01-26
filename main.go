package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v2"
)

const versionString = "zwc 1.0.1"

// Stats contains statistics about a single file, such as the number of lines
type Stats struct {
	byteCounter   uint64
	runeCounter   uint64
	wordCounter   uint64
	lineCounter   uint64
	maxLineLength uint64
}

// Examine collects statistics about a single gzipped file
func Examine(filename string) (*Stats, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(data)
	decompressorReader, err := gzip.NewReader(buf)
	if err != nil {
		return nil, fmt.Errorf("zwc: %s: not in gzip format", filename)
	}

	decompressedBytes, err := io.ReadAll(decompressorReader)
	decompressorReader.Close()
	if err != nil {
		return nil, fmt.Errorf("zwc: %s: could not decompress", filename)
	}

	var (
		stats                 Stats
		runesSinceLastNewline uint64
		inWord                bool
	)

	stats.byteCounter = uint64(len(decompressedBytes))

	for _, r := range string(decompressedBytes) {
		runesSinceLastNewline++
		stats.runeCounter++
		switch r {
		case ' ':
			if inWord {
				stats.wordCounter++
				inWord = false
			}
		case '\n':
			stats.lineCounter++
			runesSinceLastNewline--
			if runesSinceLastNewline > stats.maxLineLength {
				stats.maxLineLength = runesSinceLastNewline
			}
			runesSinceLastNewline = 0
			if inWord {
				stats.wordCounter++
				inWord = false
			}
		default:
			if !inWord {
				inWord = true
			}
		}
	}
	if inWord {
		stats.wordCounter++
	}

	return &stats, nil
}

func main() {
	if appErr := (&cli.App{
		Name:  "zwc",
		Usage: "count lines, words, bytes and runes in gzipped text files",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "version", Aliases: []string{"V"}},
			&cli.BoolFlag{Name: "lines", Aliases: []string{"l"}},
			&cli.BoolFlag{Name: "chars", Aliases: []string{"m", "r", "runes"}},
			&cli.BoolFlag{Name: "bytes", Aliases: []string{"c"}},
			&cli.BoolFlag{Name: "max-line-length", Aliases: []string{"L"}},
			&cli.BoolFlag{Name: "words", Aliases: []string{"w"}},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("version") {
				fmt.Println(versionString)
				os.Exit(0)
			}
			filenames := []string{}
			// Check if any arguments are given
			if c.NArg() > 0 {
				filenames = c.Args().Slice()
			}
			formatString := "%d %s\n"
			if len(filenames) > 0 {
				formatString = "%6d %s\n"
			}
			for _, filename := range filenames {
				stats, err := Examine(filename)
				if err != nil {
					return err
				}
				if c.Bool("lines") {
					fmt.Printf(formatString, stats.lineCounter, filename)
					continue
				}
				if c.Bool("bytes") {
					fmt.Printf(formatString, stats.byteCounter, filename)
					continue
				}
				if c.Bool("chars") {
					fmt.Printf(formatString, stats.runeCounter, filename)
					continue
				}
				if c.Bool("words") {
					fmt.Printf(formatString, stats.wordCounter, filename)
					continue
				}
				if c.Bool("max-line-length") {
					fmt.Printf(formatString, stats.maxLineLength, filename)
					continue
				}
				fmt.Printf("%4d %4d %4d %s\n", stats.lineCounter, stats.wordCounter, stats.runeCounter, filename)
			}
			return nil
		},
	}).Run(os.Args); appErr != nil {
		fmt.Fprintln(os.Stderr, appErr.Error())
		os.Exit(1)
	}
}
