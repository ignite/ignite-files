package gex

import _ "embed" // embed is required for binary embedding.

//go:embed gex-v1.3.0-amd64-linux.tar.gz
var binaryCompressed []byte
