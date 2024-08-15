package plugins

import (
	"github.com/lemmego/lemmego/api/app"
	"github.com/lemmego/lemmego/internal/plugins/auth"
)

// Load plugins
func Load() app.PluginRegistry {
	registry := app.PluginRegistry{}
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
