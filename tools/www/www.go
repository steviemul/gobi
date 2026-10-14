// Package www embeds the static web content served by the gobi server.
package www

import "embed"

// Content holds the files under gobi/, addressed as "gobi/index.html" and so on.
//
//go:embed gobi
var Content embed.FS
