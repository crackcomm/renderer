package template

import "io"

// Text - Returns template interface that does not actually execute templates.
func Text(text string) Template {
	return TextBytes([]byte(text))
}

// TextBytes - Returns template interface that does not actually execute templates.
func TextBytes(text []byte) Template {
	return &textTemplate{text: text}
}

type textTemplate struct {
	text []byte
}

// Execute - Writes text to writer. Ignores context.
func (t *textTemplate) Execute(_ Context, w io.Writer) (err error) {
	_, err = w.Write(t.text)
	return
}
