package configs

import (
	"github.com/lemmego/api/config"
	"github.com/lemmego/gpa"
	"github.com/lemmego/gpagorm"
)

func init() {
	config.Set("sql.provider", func(instance ...string) gpa.SQLProvider {
		return gpa.MustGet[*gpagorm.Provider](instance...)
	})
}
