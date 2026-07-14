package migrations

import "embed"

// Files is embedded into the API binary so migrations are available in the
// distroless production image without copying a separate directory.
//
//go:embed *.sql
var Files embed.FS
