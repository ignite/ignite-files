package hermes

import _ "embed" // embed is required for binary embedding.

//go:embed hermes-v1.6.0-aarch64-unknown-linux-gnu.tar.gz
var binaryCompressed []byte
