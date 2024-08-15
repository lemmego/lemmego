package inputs

import (
	"fmt"

	"github.com/ggicci/httpin"
	"github.com/lemmego/lemmego/api"
	"github.com/lemmego/lemmego/api/db"
)

type Input struct {
	App api.AppManager
}

type UploadedFile struct {
	*httpin.File
}

func ValidateUnique(db *db.DB, table string, column string) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		var result string
		row := db.Table(table).Where(fmt.Sprintf("%s = ?", column), v).Select(column).Row()
		if row.Err() != nil {
			return false, row.Err().Error()
		}
		row.Scan(&result)
		if result == "" {
			return true, ""
		}
		return false, "Field must be unique"
	}
}
