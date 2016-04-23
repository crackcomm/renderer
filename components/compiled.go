package components

import "github.com/crackcomm/renderer/template"

// Compiled - Compiled component ready to render.
type Compiled struct {
	// Component - Source of the compiled component.
	*Component

	// Context - Compiled component context.
	Context template.Context

	// Main - Main template compiled.
	Main template.Template

	// With - `With` templates map.
	With template.Map

	// Extends - Compiled `Extends` component.
	Extends *Compiled

	// Styles - Compiled styles templates.
	Styles []template.Template

	// Scripts - Compiled scripts templates.
	Scripts []template.Template

	// Require - Compiled `Require` components.
	Require map[string]*Compiled
}
