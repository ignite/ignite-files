package hermes

import _ "embed" // embed is required for binary embedding.

//go:embed hermes-v1.6.0-aarch64-apple-darwin.tar.gz
var binaryCompressed []byte
