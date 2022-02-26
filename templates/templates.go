package templates

import (
	"html/template"
	"path/filepath"
)

const (
	templatesDirectory = "templatefiles"
	DebugTemplateFile  = "debug.html"
)

var Templates = template.Must(
	template.ParseFiles(
		filepath.Join(templatesDirectory, DebugTemplateFile),
	),
)
