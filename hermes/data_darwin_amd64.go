package hermes

import _ "embed" // embed is required for binary embedding.

//go:embed hermes-v1.7.1-x86_64-apple-darwin.tar.gz
var binaryCompressed []byte
