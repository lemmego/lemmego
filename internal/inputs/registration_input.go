package inputs

type Registration struct {
	Username string `json:"username" in:"form=username"`
	FirstName string `json:"first_name" in:"form=first_name"`
	LastName string `json:"last_name" in:"form=last_name"`
	OrgUsername string `json:"org_username" in:"form=org_username"`
}
