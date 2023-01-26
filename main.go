package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/xyproto/textoutput"
)

const versionString = "zwc 0.0.1"

type Stats struct {
	byteCounter   uint64
	runeCounter   uint64
	wordCounter   uint64
	lineCounter   uint64
	maxLineLength uint64
}

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

	var stats Stats
	stats.byteCounter = uint64(len(decompressedBytes))

	var runesSinceLastNewline uint64 = 0
	var inWord bool

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
	o := textoutput.New()
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
				o.Println(versionString)
				os.Exit(0)
			}
			filenames := []string{}
			// Check if any arguments are given
			if c.NArg() > 0 {
				filenames = c.Args().Slice()
			}
			for _, filename := range filenames {
				stats, err := Examine(filename)
				if err != nil {
					return err
				}
				if c.Bool("lines") {
					fmt.Printf("%d %s\n", stats.lineCounter, filename)
					continue
				}
				if c.Bool("bytes") {
					fmt.Printf("%d %s\n", stats.byteCounter, filename)
					continue
				}
				if c.Bool("chars") {
					fmt.Printf("%d %s\n", stats.runeCounter, filename)
					continue
				}
				if c.Bool("words") {
					fmt.Printf("%d %s\n", stats.wordCounter, filename)
					continue
				}
				if c.Bool("max-line-length") {
					fmt.Printf("%d %s\n", stats.maxLineLength, filename)
					continue
				}
				fmt.Printf("%4d %4d %4d %s\n", stats.lineCounter, stats.wordCounter, stats.runeCounter, filename)
			}
			return nil
		},
	}).Run(os.Args); appErr != nil {
		o.ErrExit(appErr.Error())
	}
}
