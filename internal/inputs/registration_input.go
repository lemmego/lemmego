package inputs

type RegistrationStoreInput struct {
	FirstName string `json:"first_name" in:"form=first_name"`
	Email string `json:"email" in:"form=email"`
	Password string `json:"password" in:"form=password"`
}


