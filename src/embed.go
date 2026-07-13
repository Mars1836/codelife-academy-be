package assets

import "embed"

// Documents contains the temporary Markdown knowledge base.
//
//go:embed documents/*.md
var Documents embed.FS
