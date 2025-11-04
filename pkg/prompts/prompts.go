//go:build !js
package prompts

import (
	"bytes"
	_ "embed"
	"text/template"
)

//go:embed system.txt
var systemFile string

//go:embed user.txt
var userFile string

var (
	tSystem = template.Must(template.New("system").Parse(systemFile))
	tUser   = template.Must(template.New("user").Parse(userFile))
)
