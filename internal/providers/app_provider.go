package providers

import (
	"fmt"
	"github.com/lemmego/api/app"
)

type AppProvider struct{}

func (p *AppProvider) Register(a app.AppManager) {
	fmt.Println(a.Config().Get("app"))
	//TODO implement me
}

func (p *AppProvider) Boot(a app.AppManager) {
	//TODO implement me
}
