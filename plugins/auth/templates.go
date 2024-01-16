package auth

import (
	_ "embed"
)

//go:embed templates/login.page.tmpl
var loginTmpl []byte

//go:embed templates/register.page.tmpl
var registerTmpl []byte
