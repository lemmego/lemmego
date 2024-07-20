package inputs

type Login struct {
	Email string `json:"email" in:"form=email"`
	Password string `json:"password" in:"form=password"`
	OrgUsername string `json:"org_username" in:"form=org_username"`
}
