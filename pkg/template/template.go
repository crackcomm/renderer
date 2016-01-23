package template

import (
	"bytes"
	"io"
)

// Context - Template context.
type Context map[string]interface{}

// Template - Template structure.
type Template interface {
	Execute(Context, io.Writer) error
}

// Execute - Executes given template with context and returns buffer.
func Execute(t Template, ctx Context) (buf bytes.Buffer, err error) {
	err = t.Execute(ctx, &buf)
	return
}

// ExecuteToString - Executes template with context and returns string result.
func ExecuteToString(t Template, ctx Context) (res string, err error) {
	buf, err := Execute(t, ctx)
	if err != nil {
		return
	}
	return buf.String(), nil
}

// ExecuteToBytes - Executes template with context and returns string result.
func ExecuteToBytes(t Template, ctx Context) (res []byte, err error) {
	buf, err := Execute(t, ctx)
	if err != nil {
		return
	}
	return buf.Bytes(), nil
}
