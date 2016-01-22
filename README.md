# Renderer

### TODO: web interface

Components flow from `storage` to `compiler` so `compiler` does not provide `Insert` method.

What we need is to create persistent storage interface with watching for changes.

Components persistent storage will be able to process and forward following event types:

	CreateFile // not component.json
	ChangeFile // not component.json
	RemoveFile // not component.json

	CreateComponent // component.json change
	ChangeComponent // component.json change
	RemoveComponent // component.json change


Compiler if enabled will recompile templates and store in-memory on event from persistent storage.

### TODO: others

* `file://` in `styles`, `scripts` and `main` instead of `data:...` for content

## Usage

### Components

Component definition in Go:

```Go
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
	// When local files will be read and parsed as templates.
	Styles []string

	// Scripts - List of relative paths or URLs to JS files.
	// When local files will be read and parsed as templates.
	Scripts []string

	// Require - Components required by this component.
	// Those will be rendered before and set in context under keys from map.
	Require map[string]*Component `json:"require,omitempty"`

	// Context - Base context for the component.
	Context map[string]interface{} `json:"context,omitempty"`

	// With - Like context but values should be templates.
	With map[string]string `json:"with,omitempty"`
}
```

```json
{
  "name": "example.root",
  "main": "component.html",
  "styles": [
    "https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css",
    "https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap-theme.min.css"
  ],
  "scripts": [
    "https://ajax.googleapis.com/ajax/libs/jquery/1.11.3/jquery.min.js",
    "https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/js/bootstrap.min.js"
  ]
}
```

```json
{
  "name": "movies.movie",
  "main": "component.html",
  "extends": "movies.root",
  "requires": {
    "video": {
      "name": "movies.video",
      "with": {
        "title": "{{ movie.title }} ({{ movie.year }})"
      }
    }
  }
}
```

### Rendering

```json
{
  "name": "admin.domains",
  "context": {
    "domains": [
      {"Name": "test.pl", "Links": ["test.pl"]},
      {"Name": "test2.pl", "Links": ["test.pl"]},
      {"Name": "test3.pl", "Links": ["test.pl"]}
    ]
  }
}
```
