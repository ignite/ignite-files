package gex

import _ "embed" // embed is required for binary embedding.

//go:embed gex-arm64-linux.tar.gz
var binaryCompressed []byte
