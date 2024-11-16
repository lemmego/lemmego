package providers

import "github.com/lemmego/api/app"

func init() {
	// Add your services here
	app.RegisterService(func(a app.App) error {
		// Register bindings
		// e.g.:
		// a.AddService(&SomeService)
		return nil
	})

	app.BootService(func(app app.App) error {
		// Perform any start-up related tasks with some other services
		// e.g.:
		// myService := &MyService{}
		// a.Service(myService)
		// myService.Execute()
		return nil
	})
}
