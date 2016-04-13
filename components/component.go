package components

import "bitbucket.org/moovie/util/template"

// Component - Component definition.
type Component struct {
	// Name - Name of the component as registered in global scope.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// Main - Main entrypoint of rendering the component.
	Main string `json:"main,omitempty" yaml:"main,omitempty"`

	// Extends - Parent of the component.
	// Parent will be rendered with this component html as `children` in context.
	Extends string `json:"extends,omitempty" yaml:"extends,omitempty"`

	// Styles - List of relative paths or URLs to CSS files.
	// When local files will be read and parsed as templates.
	Styles []string `json:"styles,omitempty" yaml:"styles,omitempty"`

	// Scripts - List of relative paths or URLs to JS files.
	// When local files will be read and parsed as templates.
	Scripts []string `json:"scripts,omitempty" yaml:"scripts,omitempty"`

	// Require - Components required by this component.
	// Those will be rendered before and set in context under keys from map.
	Require map[string]Component `json:"require,omitempty" yaml:"require,omitempty"`

	// Context - Base context for the component.
	Context template.Context `json:"context,omitempty" yaml:"context,omitempty"`

	// With - Like context but values should be templates.
	With map[string]string `json:"with,omitempty" yaml:"with,omitempty"`
}
