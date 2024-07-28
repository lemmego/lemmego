package plugins

import (
	"github.com/invopop/validation"
	"lemmego/api"
	"lemmego/api/validator"
	"lemmego/internal/plugins/auth"
	// "pressebo/internal/plugins/auth"
)

type RegistrationInput struct {
	OrgName              string `json:"org_name" in:"form=org_name"`
	Subdomain            string `json:"subdomain" in:"form=subdomain"`
	FirstName            string `json:"first_name" in:"form=first_name"`
	LastName             string `json:"last_name" in:"form=last_name"`
	Username             string `json:"username" in:"form=username"`
	Password             string `json:"password" in:"form=password"`
	PasswordConfirmation string `json:"password_confirmation" in:"form=password_confirmation"`
}

// Validate the input
func (r *RegistrationInput) Validate() error {
	return validation.Errors{
		"org_name": validation.Validate(
			r.OrgName,
			validation.Required,
		),
		"subdomain": validation.Validate(
			r.Subdomain,
			validation.Required,
		),
		"first_name": validation.Validate(
			r.FirstName,
			validation.Required,
		),
		"last_name": validation.Validate(
			r.LastName,
			validation.Required,
		),
		"username": validation.Validate(
			r.Username,
			validation.Required,
		),
		"password": validation.Validate(
			r.Password,
			validation.Required,
		),
		"password_confirmation": validation.Validate(
			r.PasswordConfirmation,
			validation.Required,
			validation.By(validator.StringEquals(r.Password, "Passwords do not match")),
		),
	}.Filter()
}

// Load plugins
func Load() api.PluginRegistry {
	registry := api.PluginRegistry{}
	authPlugin := auth.New()
	// authPlugin := auth.New(auth.WithUserCreator(func(c *api.Context, opts *auth.Options) (bool, validation.Errors) {
	// 	input := &RegistrationInput{}
	// 	if validated, err := c.Validate(input); err != nil {
	// 		return false, err.(validation.Errors)
	// 	} else {
	// 		input = validated.(*RegistrationInput)
	// 	}

	// 	db := opts.DB

	// 	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	// 	org, err := repositories.CreateOrg(db, repositories.OrgAttrs{
	// 		Name:      input.OrgName,
	// 		Subdomain: input.Subdomain,
	// 		Email:     input.Username,
	// 	})

	// 	if err != nil {
	// 		log.Println(err)
	// 		return false, validation.Errors{
	// 			"org": fmt.Errorf("Org was not created: %w", err),
	// 		}
	// 	}

	// 	user, err := repositories.CreateUser(db, repositories.UserAttrs{
	// 		FirstName: input.FirstName,
	// 		LastName:  input.LastName,
	// 		Email:     input.Username,
	// 		Password:  string(encryptedPassword),
	// 		OrgID:     org.ID,
	// 	})

	// 	if err != nil {
	// 		log.Println(err)
	// 		return false, validation.Errors{
	// 			"user": fmt.Errorf("User was not created: %w", err),
	// 		}
	// 	}

	// 	if user.ID == 0 {
	// 		return false, validation.Errors{
	// 			"user": errors.New("User was not created."),
	// 		}
	// 	}

	// 	return true, nil
	// }), auth.WithUserResolver(func(c *api.Context, opts *auth.Options) (*auth.AuthUser, *auth.Credentials, validation.Errors) {
	// 	return nil, nil, nil
	// }))

	registry.Add(authPlugin)
	return registry
}
