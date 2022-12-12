package types

import (
	"database/sql/driver"
	"strings"
)

type strSlice []string

func (str *strSlice) Scan(val interface{}) error {
	s := val.([]uint8)
	ss := strings.Split(string(s), "|")
	*str = ss
	return nil
}

func (str strSlice) Value() (driver.Value, error) {
	strs := strings.Join(str, "|")
	return strs, nil
}
