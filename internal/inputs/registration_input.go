package inputs

type RegistrationStoreInput struct {
	Gender string `json:"gender" in:"form=gender"`
	Email string `json:"email" in:"form=email"`
	Password string `json:"password" in:"form=password"`
}


