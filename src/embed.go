package assets

import "embed"

// Documents contains the temporary Markdown knowledge base.
//
//go:embed all:documents
var Documents embed.FS
