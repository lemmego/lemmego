package gormconnector

import (
	"fmt"
	"github.com/lemmego/api/app"
	"github.com/lemmego/api/config"
	"github.com/lemmego/gpa"
	"github.com/lemmego/gpagorm"
	"reflect"
)

type Provider struct {
	config *gpa.Config
}

func (g *Provider) WithConfig(config *gpa.Config) *Provider {
	g.config = config
	return g
}

func (g *Provider) Provide(a app.App) error {
	dbConfig := sqlConfig()
	if g.config != nil {
		dbConfig = *g.config
	}
	provider, err := gpagorm.NewProvider(dbConfig)
	if err != nil {
		panic(err)
	}

	gpa.RegisterDefault(provider)

	a.AddService(provider)
	return nil
}

func sqlConfig(connName ...string) gpa.Config {
	name := "default"
	if len(connName) > 0 && connName[0] != "" {
		name = connName[0]
	}

	defaultConnection := config.Get(fmt.Sprintf("sql.%s", name))
	connection := config.Get(fmt.Sprintf("sql.connections.%s", defaultConnection)).(config.M)
	driver := connection.String("driver")
	database := connection.String("database")

	if database == "" || driver == "" {
		panic("database: database and driver must be present")
	}

	dbConfig := gpa.Config{
		Driver:   driver,
		Database: database,
	}

	if driver != "sqlite" {
		dbConfig.Host = config.Get(fmt.Sprintf("sql.connections.%s.host", defaultConnection)).(string)
		dbConfig.Port = config.Get(fmt.Sprintf("sql.connections.%s.port", defaultConnection)).(int)
		dbConfig.Username = config.Get(fmt.Sprintf("sql.connections.%s.user", defaultConnection)).(string)
		dbConfig.Password = config.Get(fmt.Sprintf("sql.connections.%s.password", defaultConnection)).(string)
		dbConfig.Options = config.Get(fmt.Sprintf("sql.connections.%s.options", defaultConnection)).(config.M)
	}

	return dbConfig
}

func Get(a app.App) *gpagorm.Provider {
	return a.Service(reflect.TypeOf(&gpagorm.Provider{})).(*gpagorm.Provider)
}
