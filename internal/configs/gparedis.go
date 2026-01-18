package configs

import (
	"github.com/lemmego/api/config"
	"github.com/lemmego/gpa"
	"github.com/lemmego/gparedis"
)

func init() {
	config.Set("keyvalue.gpaprovider", func(instance ...string) gpa.KeyValueProvider {
		return gpa.MustGet[*gparedis.Provider](instance...)
	})
}
