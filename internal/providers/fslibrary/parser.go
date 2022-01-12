package fslibrary

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/lemoony/snipkit/internal/utils/system"
)

const (
	titleHeaderLineNumberTitle = 2
)

func getSnippetName(system *system.System, filePath string) string {
	file, err := system.Fs.Open(filePath)
	fileName := filepath.Base(filePath)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = file.Close()
	}()

	if start, end, ok := getTitleHeaderRange(file); ok {
		titleHeader := make([]byte, end-start)
		_, err := file.ReadAt(titleHeader, int64(start))
		if err != nil {
			panic(err)
		}

		if title, ok := titleHeaderToTitle(titleHeader); ok {
			return title
		}
	}

	return fileName
}

func pruneTitleHeader(r io.Reader) string {
	all, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	first, last, ok := getTitleHeaderRange(bytes.NewReader(all))
	if ok {
		return strings.TrimSpace(string(all[0:first])) + strings.TrimSpace(string(all[last:]))
	}

	return string(all)
}

func getTitleHeaderRange(r io.Reader) (int, int, bool) {
	scanner := bufio.NewScanner(r)
	scanner.Split(splitterWithCRLF)

	firstPosition := 0
	lastPosition := 0
	ok := false

	title := ""
	titleLine := 0

	lineNumber := 0
	currentPosition := 0
	lastLineLength := 0

	for scanner.Scan() {
		lineNumber++
		currentPosition += lastLineLength
		scanner.Bytes()
		if lineNumber-titleLine > maxLineNumberTitleComment {
			break
		}

		line := scanner.Text()
		lastLineLength = len(line)

		if strings.HasPrefix(line, "#") {
			switch {
			case titleLine == 0 && len(strings.TrimSpace(line)) == 1:
				titleLine++
				firstPosition = currentPosition
			case titleLine == 1:
				title = strings.TrimSpace(strings.TrimPrefix(line, "#"))
				titleLine++
			case titleLine == 2 && len(strings.TrimSpace(line)) == 1:
				if title != "" {
					ok = true
					lastPosition = currentPosition + len(line)
					break
				}
			}
		}
	}

	return firstPosition, lastPosition, ok
}

func titleHeaderToTitle(bytes []byte) (string, bool) {
	scanner := bufio.NewScanner(strings.NewReader(string(bytes)))
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		if lineNumber == titleHeaderLineNumberTitle {
			return strings.TrimSpace(strings.TrimPrefix(scanner.Text(), "#")), true
		}
	}
	return "", false
}

func splitterWithCRLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		advance := i + 1
		token := data[0 : i+1]
		return advance, token, nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
