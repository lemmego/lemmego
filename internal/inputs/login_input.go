package inputs

type LoginStoreInput struct {
  Email string `json:"email" in:"form=email"`
  Password string `json:"password" in:"form=password"`
}
