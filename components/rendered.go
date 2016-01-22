package components

// Rendered - Rendered component.
type Rendered struct {
	// Body - Main body of the rendered component.
	Body string `json:"body,omitempty"`

	// Styles - List of styles.
	// They can be urls or list of css styles with prefix "data:text/css;".
	Styles []string `json:"styles,omitempty"`

	// Scripts - List of scripts.
	// They can be urls or list of js scripts with prefix "data:text/javascript;".
	Scripts []string `json:"scripts,omitempty"`

	// Links - List of routes(/links) pointing to files.
	Links map[string]string `json:"links,omitempty"`
}
