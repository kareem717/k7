package flags

import (
	"fmt"
	"strings"
)

type DBMS string

const (
	Postgres DBMS = "postgres"
)

var AllowedDBDrivers = []string{string(Postgres)}

func (f DBMS) String() string {
	return string(f)
}

func (f *DBMS) Type() string {
	return "DBMS"
}

func (f *DBMS) Set(value string) error {
	for _, db := range AllowedDBDrivers {
		if db == value {
			*f = DBMS(value)
			return nil
		}
	}

	return fmt.Errorf("database to use. Allowed values: %s", strings.Join(AllowedDBDrivers, ", "))
}