package template

import (
	"io"

	"github.com/flosch/pongo2"
)

// FromFile - Creates a new template structure from file.
func FromFile(fname string) (t Template, err error) {
	template, err := pongo2.FromFile(fname)
	if err != nil {
		return
	}
	return &pongoTemplate{Template: template}, nil
}

// FromString - Creates a new template structure from string.
func FromString(body string) (t Template, err error) {
	template, err := pongo2.FromString(body)
	if err != nil {
		return
	}
	return &pongoTemplate{Template: template}, nil
}

// FromBytes - Creates a new template structure from byte array.
func FromBytes(body []byte) (t Template, err error) {
	return FromString(string(body))
}

type pongoTemplate struct {
	*pongo2.Template
}

// Execute - Executes template with context.
func (t *pongoTemplate) Execute(ctx Context, w io.Writer) (err error) {
	err = t.Template.ExecuteWriter(pongo2.Context(ctx), w)
	return
}
