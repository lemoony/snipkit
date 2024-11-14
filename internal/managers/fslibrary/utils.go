package fslibrary

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

func formatSnippet(script string, title string) string {
	var output bytes.Buffer
	reader := bufio.NewReader(strings.NewReader(script))

	shebangLine := ""
	firstLine, err := reader.ReadString('\n')
	if err == nil && strings.HasPrefix(firstLine, "#!/") {
		shebangLine = firstLine
	} else {
		// If the first line is not a shebang, reset the reader to include the first line in the script
		reader = bufio.NewReader(strings.NewReader(firstLine + script[len(firstLine):]))
	}

	if shebangLine != "" {
		output.WriteString(shebangLine)
		output.WriteString("\n")
	}

	snippetComment := fmt.Sprintf(`#
# %s
#
`, title)
	output.WriteString(snippetComment)
	output.WriteString("\n")

	// Write the rest of the script
	for {
		line, readErr := reader.ReadString('\n')
		output.WriteString(line)
		if readErr == io.EOF {
			break
		}
	}

	return output.String()
}
