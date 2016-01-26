package template

//
// below some code
// imported from "github.com/tdewolff/parse"
//

var whitespaceTable = [256]bool{
	// ASCII
	false, false, false, false, false, false, false, false,
	false, true, true, false, true, true, false, false, // tab, new line, form feed, carriage return
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	true, false, false, false, false, false, false, false, // space
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	// non-ASCII
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
}

// isWhitespace returns true for space, \n, \r, \t, \f.
func isWhitespace(c byte) bool {
	return whitespaceTable[c]
}

// CleanWhitespaces - Replaces character series of space \n, \t, \f, \r into a
// single space or newline (when the serie contained a \n or \r).
func CleanWhitespaces(b []byte) []byte {
	j := 0
	prevWS := false
	hasNewline := false
	for i := 0; i < len(b); i++ {
		c := b[i]
		if isWhitespace(c) {
			prevWS = true
			if c == '\n' || c == '\r' {
				hasNewline = true
			}
		} else {
			if prevWS {
				prevWS = false
				if hasNewline {
					hasNewline = false
				} else {
					b[j] = ' '
					j++
				}
			}
			b[j] = b[i]
			j++
		}
	}
	if prevWS {
		if hasNewline {
			b[j] = '\n'
		} else {
			b[j] = ' '
		}
		j++
	}
	return b[:j]
}
