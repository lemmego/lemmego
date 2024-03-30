package inputs

type LoginStoreInput struct {
  Username string `json:"username" in:"form=username"`
  Password string `json:"password" in:"form=password"`
}
