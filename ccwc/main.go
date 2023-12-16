package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type option struct {
	doCount   bool
	countFunc func() int64
}

func main() {
	l := flag.Bool("l", false, "print the newline counts")
	w := flag.Bool("w", false, "print the word counts")
	c := flag.Bool("c", false, "print the byte counts")
	m := flag.Bool("m", false, "print the character counts")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("ccwc: Not enough arguments")
		os.Exit(1)
	}
	fileName := flag.Arg(0)

	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("ccwc: %s\n", err)
		os.Exit(1)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		fmt.Printf("ccwc: %s\n", err)
		os.Exit(1)
	}
	if stat.IsDir() {
		fmt.Printf("ccwc: %s: Is a directory\n", f.Name())
	}

	none := !*l && !*w && !*c && !*m
	options := []option{
		{
			doCount: none || *l,
			countFunc: func() int64 {
				return count(f, scanNewLines)
			},
		},
		{
			doCount: none || *w,
			countFunc: func() int64 {
				return count(f, bufio.ScanWords)
			},
		},
		{
			doCount: *m,
			countFunc: func() int64 {
				return count(f, bufio.ScanRunes)
			},
		},
		{
			doCount: none || *c,
			countFunc: func() int64 {
				return stat.Size()
			},
		},
	}

	result := make([]string, 0)
	for _, opt := range options {
		if opt.doCount {
			var cnt int64
			if !stat.IsDir() {
				cnt = opt.countFunc()
			}
			result = append(result, fmt.Sprintf("%d", cnt))
		}
	}
	result = append(result, fileName)

	fmt.Printf("%s\n", strings.Join(result, " "))
	if stat.IsDir() {
		os.Exit(1)
	}
}

func count(r io.ReadSeeker, split bufio.SplitFunc) int64 {
	_, _ = r.Seek(0, io.SeekStart)

	s := bufio.NewScanner(r)
	s.Split(split)
	var c int64
	for s.Scan() {
		c++
	}
	return c
}

func scanNewLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF || len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		return i + 1, dropCR(data[0:i]), nil
	}
	return 0, nil, nil
}

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}
