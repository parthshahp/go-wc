package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Options struct {
	l bool
	w bool
	c bool
	b bool
}

type stats struct {
	lines int64
	words int64
	chars int64
	bytes int64
	filename string
}

func process(r *bufio.Reader) (*stats, error) {
	var lines, words, chars, bytes, w_state int64
//	w_mask := [256]uint8{'\n': 1, '\t': 1, ' ': 1, '\v': 1, '\r': 1, '\f': 1}
//	l_mask := [256]uint8{'\n': 1}

	for {
	    next_char, bytes_read, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				res := &stats{lines: lines, words: words + w_state, chars: chars, bytes: bytes}
				return res, nil
			}
			return &stats{}, err
		}
		bytes += int64(bytes_read)
		chars++

		if unicode.IsSpace(next_char) && w_state == 1 {
			words++
			w_state = 0
		} else if !unicode.IsSpace(next_char) {
			w_state = 1
		}

		if next_char == '\n' {
			lines++
		}
	}
}

func readInput (o *Options) error {
	r := bufio.NewReader(os.Stdin)
	s, err := process(r)
	if err != nil { return err }

	printStats(s, o)

	return nil
}

func readFiles(files []string, o *Options) error {
	for _, file := range files {

		r, err := getReader(file)
		if err != nil { return err }

		s, err := process(r)
		if err != nil { return err }

		s.filename = file

		printStats(s, o)
	}

	return nil
}

func getReader(file string) (*bufio.Reader, error) {

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	return bufio.NewReader(f), nil
}

func printStats(s *stats, o *Options) {
	statPrint := []string{}

	if o.l {
		statPrint = append(statPrint, strconv.FormatInt(s.lines, 10))
	}
	if o.w {
		statPrint = append(statPrint, strconv.FormatInt(s.words, 10))
	}
	if o.b {
		statPrint = append(statPrint, strconv.FormatInt(s.bytes, 10))
	}
	if o.c {
		statPrint = append(statPrint, strconv.FormatInt(s.chars, 10))
	}
	if s.filename != "" {
		statPrint = append(statPrint, s.filename)
	}

	fmt.Println(strings.Join(statPrint, "\t"))

}
func main() {
	var commandLineOptions Options
	flag.BoolVar(&commandLineOptions.b, "b", false, "Count bytes")
	flag.BoolVar(&commandLineOptions.c, "c", false, "Count chars")
	flag.BoolVar(&commandLineOptions.w, "w", false, "Count words")
	flag.BoolVar(&commandLineOptions.l, "l", false, "Count lines")
	flag.Parse()

	if !(commandLineOptions.b || commandLineOptions.c ||
			commandLineOptions.l || commandLineOptions.w) {
		commandLineOptions.b = true
		commandLineOptions.l = true
		commandLineOptions.w = true
	}

	files := flag.CommandLine.Args()

	var err error
	if len(files) > 0 {
		err = readFiles(files, &commandLineOptions)
	} else {
		err = readInput(&commandLineOptions)
	}
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "wc: %s\n", err)
		return
	}
	
}
