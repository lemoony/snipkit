//nolint:all // since it sticks to view styling guide
package finder

// File duplicated & adapted from https://github.com/lemoony/tview/blob/master/util.go
// All functions not used by Finder have been removed.

import (
	"regexp"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
	"github.com/rivo/uniseg"
)

// Positions of substrings in regular expressions.
const (
	colorForegroundPos = 1
	colorBackgroundPos = 3
	colorFlagPos       = 5
)

var (
	colorPattern  = regexp.MustCompile(`\[([a-zA-Z]+|#[0-9a-zA-Z]{6}|\-)?(:([a-zA-Z]+|#[0-9a-zA-Z]{6}|\-)?(:([lbidrus]+|\-)?)?)?\]`)
	regionPattern = regexp.MustCompile(`\["([a-zA-Z0-9_,;: \-\.]*)"\]`)
	escapePattern = regexp.MustCompile(`\[([a-zA-Z0-9_,;: \-\."#]+)\[(\[*)\]`)
)

// printWithStyle works like Print() but it takes a style instead of just a
// foreground color. The skipWidth parameter specifies the number of cells
// skipped at the beginning of the text. It also returns the start and end index
// (exclusively) of the text actually printed. If maintainBackground is "true",
// The existing screen background is not changed (i.e. the style's background
// color is ignored).
func printWithStyle(screen tcell.Screen, text string, x, y, skipWidth, maxWidth, align int, style tcell.Style, maintainBackground bool) (int, int, int, int) {
	totalWidth, totalHeight := screen.Size()
	if maxWidth <= 0 || len(text) == 0 || y < 0 || y >= totalHeight {
		return 0, 0, 0, 0
	}

	// Decompose the text.
	colorIndices, colors, _, _, escapeIndices, strippedText, strippedWidth := decomposeString(text, false, false)

	// We want to reduce all alignments to AlignLeft.
	if align == tview.AlignRight {
		if strippedWidth-skipWidth <= maxWidth {
			// There's enough space for the entire text.
			return printWithStyle(screen, text, x+maxWidth-strippedWidth+skipWidth, y, skipWidth, maxWidth, tview.AlignLeft, style, maintainBackground)
		}
		// Trim characters off the beginning.
		var (
			bytes, width, colorPos, escapePos, tagOffset, from, to int
			foregroundColor, backgroundColor, attributes           string
		)
		originalStyle := style
		iterateString(strippedText, func(main rune, comb []rune, textPos, textWidth, screenPos, screenWidth int) bool {
			// Update color/escape tag offset and style.
			if colorPos < len(colorIndices) && textPos+tagOffset >= colorIndices[colorPos][0] && textPos+tagOffset < colorIndices[colorPos][1] {
				foregroundColor, backgroundColor, attributes = styleFromTag(foregroundColor, backgroundColor, attributes, colors[colorPos])
				style = overlayStyle(originalStyle, foregroundColor, backgroundColor, attributes)
				tagOffset += colorIndices[colorPos][1] - colorIndices[colorPos][0]
				colorPos++
			}
			if escapePos < len(escapeIndices) && textPos+tagOffset >= escapeIndices[escapePos][0] && textPos+tagOffset < escapeIndices[escapePos][1] {
				tagOffset++
				escapePos++
			}
			if strippedWidth-screenPos <= maxWidth {
				// We chopped off enough.
				if escapePos > 0 && textPos+tagOffset-1 >= escapeIndices[escapePos-1][0] && textPos+tagOffset-1 < escapeIndices[escapePos-1][1] {
					// Unescape open escape sequences.
					escapeCharPos := escapeIndices[escapePos-1][1] - 2
					text = text[:escapeCharPos] + text[escapeCharPos+1:]
				}
				// Print and return.
				bytes, width, from, to = printWithStyle(screen, text[textPos+tagOffset:], x, y, 0, maxWidth, tview.AlignLeft, style, maintainBackground)
				from += textPos + tagOffset
				to += textPos + tagOffset
				return true
			}
			return false
		})
		return bytes, width, from, to
	} else if align == tview.AlignCenter {
		if strippedWidth-skipWidth == maxWidth {
			// Use the exact space.
			return printWithStyle(screen, text, x, y, skipWidth, maxWidth, tview.AlignLeft, style, maintainBackground)
		} else if strippedWidth-skipWidth < maxWidth {
			// We have more space than we need.
			half := (maxWidth - strippedWidth + skipWidth) / 2
			return printWithStyle(screen, text, x+half, y, skipWidth, maxWidth-half, tview.AlignLeft, style, maintainBackground)
		} else {
			// Chop off runes until we have a perfect fit.
			var choppedLeft, choppedRight, leftIndex, rightIndex int
			rightIndex = len(strippedText)
			for rightIndex-1 > leftIndex && strippedWidth-skipWidth-choppedLeft-choppedRight > maxWidth {
				if skipWidth > 0 || choppedLeft < choppedRight {
					// Iterate on the left by one character.
					iterateString(strippedText[leftIndex:], func(main rune, comb []rune, textPos, textWidth, screenPos, screenWidth int) bool {
						if skipWidth > 0 {
							skipWidth -= screenWidth
							strippedWidth -= screenWidth
						} else {
							choppedLeft += screenWidth
						}
						leftIndex += textWidth
						return true
					})
				} else {
					// Iterate on the right by one character.
					iterateStringReverse(strippedText[leftIndex:rightIndex], func(main rune, comb []rune, textPos, textWidth, screenPos, screenWidth int) bool {
						choppedRight += screenWidth
						rightIndex -= textWidth
						return true
					})
				}
			}

			// Add tag offsets and determine start style.
			var (
				colorPos, escapePos, tagOffset               int
				foregroundColor, backgroundColor, attributes string
			)
			originalStyle := style
			for index := range strippedText {
				// We only need the offset of the left index.
				if index > leftIndex {
					// We're done.
					if escapePos > 0 && leftIndex+tagOffset-1 >= escapeIndices[escapePos-1][0] && leftIndex+tagOffset-1 < escapeIndices[escapePos-1][1] {
						// Unescape open escape sequences.
						escapeCharPos := escapeIndices[escapePos-1][1] - 2
						text = text[:escapeCharPos] + text[escapeCharPos+1:]
					}
					break
				}

				// Update color/escape tag offset.
				if colorPos < len(colorIndices) && index+tagOffset >= colorIndices[colorPos][0] && index+tagOffset < colorIndices[colorPos][1] {
					if index <= leftIndex {
						foregroundColor, backgroundColor, attributes = styleFromTag(foregroundColor, backgroundColor, attributes, colors[colorPos])
						style = overlayStyle(originalStyle, foregroundColor, backgroundColor, attributes)
					}
					tagOffset += colorIndices[colorPos][1] - colorIndices[colorPos][0]
					colorPos++
				}
				if escapePos < len(escapeIndices) && index+tagOffset >= escapeIndices[escapePos][0] && index+tagOffset < escapeIndices[escapePos][1] {
					tagOffset++
					escapePos++
				}
			}
			bytes, width, from, to := printWithStyle(screen, text[leftIndex+tagOffset:], x, y, 0, maxWidth, tview.AlignLeft, style, maintainBackground)
			from += leftIndex + tagOffset
			to += leftIndex + tagOffset
			return bytes, width, from, to
		}
	}

	// Draw text.
	var (
		drawn, drawnWidth, colorPos, escapePos, tagOffset, from, to int
		foregroundColor, backgroundColor, attributes                string
	)
	iterateString(strippedText, func(main rune, comb []rune, textPos, length, screenPos, screenWidth int) bool {
		// Skip character if necessary.
		if skipWidth > 0 {
			skipWidth -= screenWidth
			from = textPos + length
			to = from
			return false
		}

		// Only continue if there is still space.
		if drawnWidth+screenWidth > maxWidth || x+drawnWidth >= totalWidth {
			return true
		}

		// Handle color tags.
		for colorPos < len(colorIndices) && textPos+tagOffset >= colorIndices[colorPos][0] && textPos+tagOffset < colorIndices[colorPos][1] {
			foregroundColor, backgroundColor, attributes = styleFromTag(foregroundColor, backgroundColor, attributes, colors[colorPos])
			tagOffset += colorIndices[colorPos][1] - colorIndices[colorPos][0]
			colorPos++
		}

		// Handle escape tags.
		if escapePos < len(escapeIndices) && textPos+tagOffset >= escapeIndices[escapePos][0] && textPos+tagOffset < escapeIndices[escapePos][1] {
			if textPos+tagOffset == escapeIndices[escapePos][1]-2 {
				tagOffset++
				escapePos++
			}
		}

		// Memorize positions.
		to = textPos + length

		// Print the rune sequence.
		finalX := x + drawnWidth
		finalStyle := style
		if maintainBackground {
			_, _, existingStyle, _ := screen.GetContent(finalX, y)
			_, background, _ := existingStyle.Decompose()
			finalStyle = finalStyle.Background(background)
		}
		finalStyle = overlayStyle(finalStyle, foregroundColor, backgroundColor, attributes)
		for offset := screenWidth - 1; offset >= 0; offset-- {
			// To avoid undesired effects, we populate all cells.
			if offset == 0 {
				screen.SetContent(finalX+offset, y, main, comb, finalStyle)
			} else {
				screen.SetContent(finalX+offset, y, ' ', nil, finalStyle)
			}
		}

		// Advance.
		drawn += length
		drawnWidth += screenWidth

		return false
	})

	return drawn + tagOffset + len(escapeIndices), drawnWidth, from, to
}

// printWithStyle works like Print() but it takes a style instead of just a
// foreground color. The skipWidth parameter specifies the number of cells
// skipped at the beginning of the text. It also returns the start and end index
// (exclusively) of the text actually printed. If maintainBackground is "true",
// The existing screen background is not changed (i.e. the style's background
// color is ignored).
func printWithStyleStripColors(screen tcell.Screen, text string, x, y, skipWidth, maxWidth, align int, style tcell.Style, stripColors bool, maintainBackground bool) (int, int, int, int) {
	totalWidth, totalHeight := screen.Size()
	if maxWidth <= 0 || len(text) == 0 || y < 0 || y >= totalHeight {
		return 0, 0, 0, 0
	}

	// Decompose the text.
	colorIndices, colors, _, _, escapeIndices, strippedText, strippedWidth := decomposeString(text, stripColors, false)

	// We want to reduce all alignments to AlignLeft.
	if align == tview.AlignRight {
		if strippedWidth-skipWidth <= maxWidth {
			// There's enough space for the entire text.
			return printWithStyleStripColors(screen, text, x+maxWidth-strippedWidth+skipWidth, y, skipWidth, maxWidth, tview.AlignLeft, style, stripColors, maintainBackground)
		}
		// Trim characters off the beginning.
		var (
			bytes, width, colorPos, escapePos, tagOffset, from, to int
			foregroundColor, backgroundColor, attributes           string
		)
		originalStyle := style
		iterateString(strippedText, func(main rune, comb []rune, textPos, textWidth, screenPos, screenWidth int) bool {
			// Update color/escape tag offset and style.
			if colorPos < len(colorIndices) && textPos+tagOffset >= colorIndices[colorPos][0] && textPos+tagOffset < colorIndices[colorPos][1] {
				foregroundColor, backgroundColor, attributes = styleFromTag(foregroundColor, backgroundColor, attributes, colors[colorPos])
				style = overlayStyle(originalStyle, foregroundColor, backgroundColor, attributes)
				tagOffset += colorIndices[colorPos][1] - colorIndices[colorPos][0]
				colorPos++
			}
			if escapePos < len(escapeIndices) && textPos+tagOffset >= escapeIndices[escapePos][0] && textPos+tagOffset < escapeIndices[escapePos][1] {
				tagOffset++
				escapePos++
			}
			if strippedWidth-screenPos <= maxWidth {
				// We chopped off enough.
				if escapePos > 0 && textPos+tagOffset-1 >= escapeIndices[escapePos-1][0] && textPos+tagOffset-1 < escapeIndices[escapePos-1][1] {
					// Unescape open escape sequences.
					escapeCharPos := escapeIndices[escapePos-1][1] - 2
					text = text[:escapeCharPos] + text[escapeCharPos+1:]
				}
				// Print and return.
				bytes, width, from, to = printWithStyleStripColors(screen, text[textPos+tagOffset:], x, y, 0, maxWidth, tview.AlignLeft, style, false, maintainBackground)
				from += textPos + tagOffset
				to += textPos + tagOffset
				return true
			}
			return false
		})
		return bytes, width, from, to
	} else if align == tview.AlignCenter {
		if strippedWidth-skipWidth == maxWidth {
			// Use the exact space.
			return printWithStyleStripColors(screen, text, x, y, skipWidth, maxWidth, tview.AlignLeft, style, false, maintainBackground)
		} else if strippedWidth-skipWidth < maxWidth {
			// We have more space than we need.
			half := (maxWidth - strippedWidth + skipWidth) / 2
			return printWithStyleStripColors(screen, text, x+half, y, skipWidth, maxWidth-half, tview.AlignLeft, style, false, maintainBackground)
		} else {
			// Chop off runes until we have a perfect fit.
			var choppedLeft, choppedRight, leftIndex, rightIndex int
			rightIndex = len(strippedText)
			for rightIndex-1 > leftIndex && strippedWidth-skipWidth-choppedLeft-choppedRight > maxWidth {
				if skipWidth > 0 || choppedLeft < choppedRight {
					// Iterate on the left by one character.
					iterateString(strippedText[leftIndex:], func(main rune, comb []rune, textPos, textWidth, screenPos, screenWidth int) bool {
						if skipWidth > 0 {
							skipWidth -= screenWidth
							strippedWidth -= screenWidth
						} else {
							choppedLeft += screenWidth
						}
						leftIndex += textWidth
						return true
					})
				} else {
					// Iterate on the right by one character.
					iterateStringReverse(strippedText[leftIndex:rightIndex], func(main rune, comb []rune, textPos, textWidth, screenPos, screenWidth int) bool {
						choppedRight += screenWidth
						rightIndex -= textWidth
						return true
					})
				}
			}

			// Add tag offsets and determine start style.
			var (
				colorPos, escapePos, tagOffset               int
				foregroundColor, backgroundColor, attributes string
			)
			originalStyle := style
			for index := range strippedText {
				// We only need the offset of the left index.
				if index > leftIndex {
					// We're done.
					if escapePos > 0 && leftIndex+tagOffset-1 >= escapeIndices[escapePos-1][0] && leftIndex+tagOffset-1 < escapeIndices[escapePos-1][1] {
						// Unescape open escape sequences.
						escapeCharPos := escapeIndices[escapePos-1][1] - 2
						text = text[:escapeCharPos] + text[escapeCharPos+1:]
					}
					break
				}

				// Update color/escape tag offset.
				if colorPos < len(colorIndices) && index+tagOffset >= colorIndices[colorPos][0] && index+tagOffset < colorIndices[colorPos][1] {
					if index <= leftIndex {
						foregroundColor, backgroundColor, attributes = styleFromTag(foregroundColor, backgroundColor, attributes, colors[colorPos])
						style = overlayStyle(originalStyle, foregroundColor, backgroundColor, attributes)
					}
					tagOffset += colorIndices[colorPos][1] - colorIndices[colorPos][0]
					colorPos++
				}
				if escapePos < len(escapeIndices) && index+tagOffset >= escapeIndices[escapePos][0] && index+tagOffset < escapeIndices[escapePos][1] {
					tagOffset++
					escapePos++
				}
			}
			bytes, width, from, to := printWithStyleStripColors(screen, text[leftIndex+tagOffset:], x, y, 0, maxWidth, tview.AlignLeft, style, stripColors, maintainBackground)
			from += leftIndex + tagOffset
			to += leftIndex + tagOffset
			return bytes, width, from, to
		}
	}

	// Draw text.
	var (
		drawn, drawnWidth, colorPos, escapePos, tagOffset, from, to int
		foregroundColor, backgroundColor, attributes                string
	)
	iterateString(strippedText, func(main rune, comb []rune, textPos, length, screenPos, screenWidth int) bool {
		// Skip character if necessary.
		if skipWidth > 0 {
			skipWidth -= screenWidth
			from = textPos + length
			to = from
			return false
		}

		// Only continue if there is still space.
		if drawnWidth+screenWidth > maxWidth || x+drawnWidth >= totalWidth {
			return true
		}

		// Handle color tags.
		for colorPos < len(colorIndices) && textPos+tagOffset >= colorIndices[colorPos][0] && textPos+tagOffset < colorIndices[colorPos][1] {
			foregroundColor, backgroundColor, attributes = styleFromTag(foregroundColor, backgroundColor, attributes, colors[colorPos])
			tagOffset += colorIndices[colorPos][1] - colorIndices[colorPos][0]
			colorPos++
		}

		// Handle escape tags.
		if escapePos < len(escapeIndices) && textPos+tagOffset >= escapeIndices[escapePos][0] && textPos+tagOffset < escapeIndices[escapePos][1] {
			if textPos+tagOffset == escapeIndices[escapePos][1]-2 {
				tagOffset++
				escapePos++
			}
		}

		// Memorize positions.
		to = textPos + length

		// Print the rune sequence.
		finalX := x + drawnWidth
		finalStyle := style
		if maintainBackground {
			_, _, existingStyle, _ := screen.GetContent(finalX, y)
			_, background, _ := existingStyle.Decompose()
			finalStyle = finalStyle.Background(background)
		}
		finalStyle = overlayStyle(finalStyle, foregroundColor, backgroundColor, attributes)
		for offset := screenWidth - 1; offset >= 0; offset-- {
			// To avoid undesired effects, we populate all cells.
			if offset == 0 {
				screen.SetContent(finalX+offset, y, main, comb, finalStyle)
			} else {
				screen.SetContent(finalX+offset, y, ' ', nil, finalStyle)
			}
		}

		// Advance.
		drawn += length
		drawnWidth += screenWidth

		return false
	})

	return drawn + tagOffset + len(escapeIndices), drawnWidth, from, to
}

// decomposeString returns information about a string which may contain color
// tags or region tags, depending on which ones are requested to be found. It
// returns the indices of the color tags (as returned by
// re.FindAllStringIndex()), the color tags themselves (as returned by
// re.FindAllStringSubmatch()), the indices of region tags and the region tags
// themselves, the indices of an escaped tags (only if at least color tags or
// region tags are requested), the string stripped by any tags and escaped, and
// the screen width of the stripped string.
func decomposeString(text string, findColors, findRegions bool) (colorIndices [][]int, colors [][]string, regionIndices [][]int, regions [][]string, escapeIndices [][]int, stripped string, width int) {
	// Shortcut for the trivial case.
	if !findColors && !findRegions {
		return nil, nil, nil, nil, nil, text, stringWidth(text)
	}

	// Get positions of any tags.
	if findColors {
		colorIndices = colorPattern.FindAllStringIndex(text, -1)
		colors = colorPattern.FindAllStringSubmatch(text, -1)
	}
	if findRegions {
		regionIndices = regionPattern.FindAllStringIndex(text, -1)
		regions = regionPattern.FindAllStringSubmatch(text, -1)
	}
	escapeIndices = escapePattern.FindAllStringIndex(text, -1)

	// Because the color pattern detects empty tags, we need to filter them out.
	for i := len(colorIndices) - 1; i >= 0; i-- {
		if colorIndices[i][1]-colorIndices[i][0] == 2 {
			colorIndices = append(colorIndices[:i], colorIndices[i+1:]...)
			colors = append(colors[:i], colors[i+1:]...)
		}
	}

	// Make a (sorted) list of all tags.
	allIndices := make([][3]int, 0, len(colorIndices)+len(regionIndices)+len(escapeIndices))
	for indexType, index := range [][][]int{colorIndices, regionIndices, escapeIndices} {
		for _, tag := range index {
			allIndices = append(allIndices, [3]int{tag[0], tag[1], indexType})
		}
	}
	sort.Slice(allIndices, func(i int, j int) bool {
		return allIndices[i][0] < allIndices[j][0]
	})

	// Remove the tags from the original string.
	var from int
	buf := make([]byte, 0, len(text))
	for _, indices := range allIndices {
		if indices[2] == 2 { // Escape sequences are not simply removed.
			buf = append(buf, []byte(text[from:indices[1]-2])...)
			buf = append(buf, ']')
			from = indices[1]
		} else {
			buf = append(buf, []byte(text[from:indices[0]])...)
			from = indices[1]
		}
	}
	buf = append(buf, text[from:]...)
	stripped = string(buf)

	// Get the width of the stripped string.
	width = stringWidth(stripped)

	return
}

// stringWidth returns the number of horizontal cells needed to print the given
// text. It splits the text into its grapheme clusters, calculates each
// cluster's width, and adds them up to a total.
func stringWidth(text string) (width int) {
	g := uniseg.NewGraphemes(text)
	for g.Next() {
		var chWidth int
		for _, r := range g.Runes() {
			chWidth = runewidth.RuneWidth(r)
			if chWidth > 0 {
				break // Our best guess at this point is to use the width of the first non-zero-width rune.
			}
		}
		width += chWidth
	}
	return
}

// iterateString iterates through the given string one printed character at a
// time. For each such character, the callback function is called with the
// Unicode code points of the character (the first rune and any combining runes
// which may be nil if there aren't any), the starting position (in bytes)
// within the original string, its length in bytes, the screen position of the
// character, and the screen width of it. The iteration stops if the callback
// returns true. This function returns true if the iteration was stopped before
// the last character.
func iterateString(text string, callback func(main rune, comb []rune, textPos, textWidth, screenPos, screenWidth int) bool) bool {
	var screenPos int

	gr := uniseg.NewGraphemes(text)
	for gr.Next() {
		r := gr.Runes()
		from, to := gr.Positions()
		width := stringWidth(gr.Str())
		var comb []rune
		if len(r) > 1 {
			comb = r[1:]
		}

		if callback(r[0], comb, from, to-from, screenPos, width) {
			return true
		}

		screenPos += width
	}

	return false
}

// styleFromTag takes the given style, defined by a foreground color (fgColor),
// a background color (bgColor), and style attributes, and modifies it based on
// the substrings (tagSubstrings) extracted by the regular expression for color
// tags. The new colors and attributes are returned where empty strings mean
// "don't modify" and a dash ("-") means "reset to default".
func styleFromTag(fgColor, bgColor, attributes string, tagSubstrings []string) (newFgColor, newBgColor, newAttributes string) {
	if tagSubstrings[colorForegroundPos] != "" {
		color := tagSubstrings[colorForegroundPos]
		if color == "-" {
			fgColor = "-"
		} else if color != "" {
			fgColor = color
		}
	}

	if tagSubstrings[colorBackgroundPos-1] != "" {
		color := tagSubstrings[colorBackgroundPos]
		if color == "-" {
			bgColor = "-"
		} else if color != "" {
			bgColor = color
		}
	}

	if tagSubstrings[colorFlagPos-1] != "" {
		flags := tagSubstrings[colorFlagPos]
		if flags == "-" {
			attributes = "-"
		} else if flags != "" {
			attributes = flags
		}
	}

	return fgColor, bgColor, attributes
}

// overlayStyle calculates a new style based on "style" and applying tag-based
// colors/attributes to it (see also styleFromTag()).
func overlayStyle(style tcell.Style, fgColor, bgColor, attributes string) tcell.Style {
	_, _, defAttr := style.Decompose()

	if fgColor != "" && fgColor != "-" {
		style = style.Foreground(tcell.GetColor(fgColor))
	}

	if bgColor != "" && bgColor != "-" {
		style = style.Background(tcell.GetColor(bgColor))
	}

	if attributes == "-" {
		style = style.Bold(defAttr&tcell.AttrBold > 0).
			Italic(defAttr&tcell.AttrItalic > 0).
			Blink(defAttr&tcell.AttrBlink > 0).
			Reverse(defAttr&tcell.AttrReverse > 0).
			Underline(defAttr&tcell.AttrUnderline > 0).
			Dim(defAttr&tcell.AttrDim > 0)
	} else if attributes != "" {
		style = style.Normal()
		for _, flag := range attributes {
			switch flag {
			case 'l':
				style = style.Blink(true)
			case 'b':
				style = style.Bold(true)
			case 'i':
				style = style.Italic(true)
			case 'd':
				style = style.Dim(true)
			case 'r':
				style = style.Reverse(true)
			case 'u':
				style = style.Underline(true)
			case 's':
				style = style.StrikeThrough(true)
			}
		}
	}

	return style
}

// iterateStringReverse iterates through the given string in reverse, starting
// from the end of the string, one printed character at a time. For each such
// character, the callback function is called with the Unicode code points of
// the character (the first rune and any combining runes which may be nil if
// there aren't any), the starting position (in bytes) within the original
// string, its length in bytes, the screen position of the character, and the
// screen width of it. The iteration stops if the callback returns true. This
// function returns true if the iteration was stopped before the last character.
func iterateStringReverse(text string, callback func(main rune, comb []rune, textPos, textWidth, screenPos, screenWidth int) bool) bool {
	type cluster struct {
		main                                       rune
		comb                                       []rune
		textPos, textWidth, screenPos, screenWidth int
	}

	// Create the grapheme clusters.
	var clusters []cluster
	iterateString(text, func(main rune, comb []rune, textPos int, textWidth int, screenPos int, screenWidth int) bool {
		clusters = append(clusters, cluster{
			main:        main,
			comb:        comb,
			textPos:     textPos,
			textWidth:   textWidth,
			screenPos:   screenPos,
			screenWidth: screenWidth,
		})
		return false
	})

	// Iterate in reverse.
	for index := len(clusters) - 1; index >= 0; index-- {
		if callback(
			clusters[index].main,
			clusters[index].comb,
			clusters[index].textPos,
			clusters[index].textWidth,
			clusters[index].screenPos,
			clusters[index].screenWidth,
		) {
			return true
		}
	}

	return false
}

func minInt(first, second int) int {
	if first < second {
		return first
	}
	return second
}
