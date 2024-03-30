package tpl

import "embed"

//go:embed server/*.tpl
var ServerTemplateFS embed.FS
