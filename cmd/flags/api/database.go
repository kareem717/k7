package flags

import (
	"fmt"
	"strings"
)

type Database string

const (
	Postgres Database = "postgres"
	None     Database = "none"
)

var AllowedDBDrivers = []string{string(Postgres), string(None)}

func (f Database) String() string {
	return string(f)
}

func (f *Database) Type() string {
	return "Database"
}

func (f *Database) Set(value string) error {
	for _, database := range AllowedDBDrivers {
		if database == value {
			*f = Database(value)
			return nil
		}
	}

	return fmt.Errorf("database to use. Allowed values: %s", strings.Join(AllowedDBDrivers, ", "))
}