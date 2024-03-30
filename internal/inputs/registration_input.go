package inputs

type RegistrationStoreInput struct {
	FirstName string `json:"first_name" in:"form=first_name"`
	LastName string `json:"last_name" in:"form=last_name"`
	Email string `json:"email" in:"form=email"`
	Username string `json:"username" in:"form=username"`
	Password string `json:"password" in:"form=password"`
}


