package main

import "embed"

//go:embed aapt.exe apksigner.jar zipalign.exe key.jks
var embeddedFiles embed.FS
