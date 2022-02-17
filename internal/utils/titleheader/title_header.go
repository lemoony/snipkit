package titleheader

import (
	"bufio"
	"bytes"
	"strings"
)

const (
	titleHeaderLineNumberTitle = 2
	maxLineNumberTitleComment  = 3
)

func ParseTitleFromHeader(content string) (string, bool) {
	if start, end, ok := getTitleHeaderRange(content); ok {
		titleHeader := content[start:end]
		return titleHeaderToTitle(titleHeader)
	}
	return "", false
}

func PruneTitleHeader(content string) string {
	first, last, ok := getTitleHeaderRange(content)
	if ok {
		result := strings.TrimSpace(content[0:first])
		if end := strings.TrimSpace(content[last:]); end != "" {
			if result != "" {
				result += "\n\n"
			}
			result += end
		}
		return result
	}

	return content
}

func getTitleHeaderRange(content string) (int, int, bool) {
	scanner := bufio.NewScanner(strings.NewReader(content))
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

func titleHeaderToTitle(content string) (string, bool) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		if lineNumber == titleHeaderLineNumberTitle {
			return strings.TrimSpace(strings.TrimPrefix(scanner.Text(), "#")), true
		}
	}
	// should never be reached since function is only called for valid headers
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
