package db

import (
	"errors"
	"regexp"
	"slices"
	"strings"
)

const (
	DialectSQLite   = "sqlite"
	DialectMySQL    = "mysql"
	DialectPostgres = "postgres"
)

var (
	supportedDialects     = []string{DialectSQLite, DialectMySQL, DialectPostgres /*, "mssql"*/}
	ErrUnsupportedDialect = errors.New("unsupported dialect")
)

// DataSource holds the necessary fields for a DSN (data source name)
type DataSource struct {
	Dialect  string
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	Params   string
}

// String returns the string representation of the data source
func (ds *DataSource) String() (string, error) {
	dialect := strings.ToLower(ds.Dialect)

	if ds.Dialect == "" {
		return "", errors.New("Dialect is required")
	}

	if !slices.Contains(supportedDialects, dialect) {
		return "", ErrUnsupportedDialect
	}

	if dialect != DialectSQLite && ds.Host == "" {
		return "", errors.New("DB Host is required")
	}

	if dialect != DialectSQLite && ds.Username == "" {
		return "", errors.New("DB Username is required")
	}

	if ds.Name == "" {
		return "", errors.New("DB Name is required")
	}

	if ds.Dialect == DialectMySQL && ds.Port == "" {
		ds.Port = "3306"
	}

	if ds.Dialect == DialectPostgres && ds.Port == "" {
		ds.Port = "5432"
	}

	// if d.Dialect == "mssql" && d.Port == "" {
	// 	d.Port = "1433"
	// }

	if ds.Dialect == DialectSQLite {
		ds.Host = ds.Name
	}

	ds.validateParams(ds.Params)

	if ds.Dialect == DialectMySQL /*|| d.Dialect == "mssql"*/ {
		ds.Params = "?" + ds.Params
	}

	if ds.Dialect == DialectPostgres {
		split := strings.Split(ds.Params, "&")
		ds.Params = " " + strings.Join(split, " ")
	}

	switch ds.Dialect {
	case DialectSQLite:
		return ds.getSqliteDSN(), nil
	case DialectMySQL:
		return ds.getMysqlDSN(), nil
	case DialectPostgres:
		return ds.getPostgresDSN(), nil
	// case "mssql":
	// 	return dsn.getMssqlDSN(), nil
	default:
		return "", ErrUnsupportedDialect
	}
}

// Make sure the params field conforms to this format: param1=value1&paramN=valueN
func (d *DataSource) validateParams(params string) error {
	if regexp.MustCompile(`^(?:[a-zA-Z0-9]+=[a-zA-Z0-9]+)(?:&[a-zA-Z0-9]+=[a-zA-Z0-9]+)*$`).MatchString(params) {
		return nil
	}

	return errors.New("invalid params format")
}

// Example: file::memory:?cache=shared
func (d *DataSource) getSqliteDSN() string {
	return "file:" + d.Name + "?" + d.Params
}

// Example: root:password@tcp(localhost:3306)/test?parseTime=true
func (d *DataSource) getMysqlDSN() string {
	return d.Username + ":" + d.Password + "@tcp(" + d.Host + ":" + string(d.Port) + ")/" + d.Name + d.Params
}

// Example: host=localhost port=5432 user=root password=password dbname=test sslmode=disable
func (d *DataSource) getPostgresDSN() string {
	hostStr := ""
	portStr := ""
	userStr := ""
	passStr := ""
	dbStr := ""
	paramsStr := ""

	if d.Host != "" {
		hostStr = "host=" + d.Host
	}

	if d.Port != "" {
		portStr = " port=" + d.Port
	}

	if d.Username != "" {
		userStr = " user=" + d.Username
	}

	if d.Password != "" {
		passStr = " password=" + d.Password
	}

	if d.Name != "" {
		dbStr = " dbname=" + d.Name
	}

	if d.Params != "" {
		paramsStr = d.Params
	}

	return hostStr + portStr + userStr + passStr + dbStr + paramsStr
}

// func (d *DSN) getMssqlDSN() string {
// 	return "sqlserver://" + d.username + ":" + d.password + "@" + d.host + ":" + string(d.port) + "?database=" + d.name
// }
