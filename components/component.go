package components

// Component - Component definition.
type Component struct {
	// Name - Name of the component as registered in global scope.
	Name string `json:"name,omitempty"`

	// Main - Main entrypoint of rendering the component.
	Main string `json:"main,omitempty"`

	// Extends - Parent of the component.
	// Parent will be rendered with this component html as `children` in context.
	Extends string `json:"extends,omitempty"`

	// Styles - List of relative paths or URLs to CSS files.
	// When local files will be read and parsed as templates.
	Styles []string `json:"styles,omitempty"`

	// Scripts - List of relative paths or URLs to JS files.
	// When local files will be read and parsed as templates.
	Scripts []string `json:"scripts,omitempty"`

	// Require - Components required by this component.
	// Those will be rendered before and set in context under keys from map.
	Require map[string]*Component `json:"require,omitempty"`

	// Context - Base context for the component.
	Context map[string]interface{} `json:"context,omitempty"`

	// With - Like context but values should be templates.
	With map[string]string `json:"with,omitempty"`
}

// Rendered - Rendered component.
type Rendered struct {
	// Body - Body of the rendered component.
	Body string `json:"body,omitempty"`

	// Links - List of routes(/links) pointing to files.
	Links map[string]string `json:"links,omitempty"`

	// Styles - List of styles.
	// They can be urls or list of css styles with prefix "data:text/css;".
	Styles []string `json:"styles,omitempty"`

	// Scripts - List of scripts.
	// They can be urls or list of js scripts with prefix "data:text/javascript;".
	Scripts []string `json:"scripts,omitempty"`
}
