package toc

import (
	"bufio"
	"bytes"
	"regexp"
)

var (
	rTocStart = regexp.MustCompile("^<!\\-\\-\\ ToC\\ start\\ \\-\\->")
	rTocEnd   = regexp.MustCompile("^<!\\-\\-\\ ToC\\ end\\ \\-\\->")
)

// Replace is taking a markdown file as input and replaces the existing table of
// contents with a new one. If there is no previous ToC, the new one will be
// injected on top of everything else.
func Replace(d []byte, toc []string) []string {
	var lines []string
	tocBlock := false
	tocFound := false
	s := bufio.NewScanner(bytes.NewReader(d))
	for s.Scan() {
		if rTocStart.Match(s.Bytes()) {
			tocBlock = true
			tocFound = true
			lines = append(lines, toc...)
		}

		if !tocBlock {
			lines = append(lines, s.Text())
		}

		if rTocEnd.Match(s.Bytes()) {
			tocBlock = false
		}
	}

	// Insert the ToC on top if no ToC tags were found in the file.
	if !tocFound {
		lines = append(toc, lines...)
	}

	return lines
}
