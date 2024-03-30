package create

import "html/template"

const GitKeep = ".gitkeep"

type Create struct {
	ProjectName string
	GoVersion   string
	Lt          template.HTML
}
