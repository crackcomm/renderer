package components

// Component - Component definition.
type Component struct {
	// Name - Name of the component as registered in global scope.
	Name string

	// Main - Main entrypoint of rendering the component.
	Main string

	// Extends - Parent of the component.
	// Parent will be rendered with this component html as `children` in context.
	Extends string

	// Styles - List of relative paths or URLs to CSS files.
	Styles []string

	// Scripts - List of relative paths or URLs to JS files.
	Scripts []string

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
	Body string

	// Links - List of routes(/links) pointing to files.
	Links map[string]string
}
