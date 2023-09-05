package conf

import "embed"

//go:embed credentials.json
var TokenData embed.FS
