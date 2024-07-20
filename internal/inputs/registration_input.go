package inputs

import (
	"fmt"
	"lemmego/api/vee"
)

type Registration struct {
	FirstName            string   `json:"first_name" in:"form=first_name"`
	LastName             string   `json:"last_name" in:"form=last_name"`
	OrgUsername          string   `json:"org_username" in:"form=org_username"`
	OrgName              string   `json:"org_name" in:"form=org_name"`
	Password             string   `json:"password" in:"form=password"`
	PasswordConfirmation string   `json:"password_confirmation" in:"form=password_confirmation"`
	FavoriteFruits       []string `json:"favorite_fruits"`
}

func (r *Registration) Validate() error {
	v := vee.New()
	v.Required("first_name", r.FirstName)
	v.Required("last_name", r.LastName)
	v.ForEach("favorite_fruits", r.FavoriteFruits, func(field string, value interface{}, index int) bool {
		return v.Contains(fmt.Sprintf("favorite_fruits.%d", index), value.(string), "a")
	})

	return v.Errors
}
