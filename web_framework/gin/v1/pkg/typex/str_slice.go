package typex

import (
	"database/sql/driver"
	"strings"
)

// 给gorm提供数组类型的支持
// https://blog.csdn.net/js010111/article/details/126076320
type StrSlice []string

func (str *StrSlice) Scan(val interface{}) error {
	s := val.([]uint8)
	ss := strings.Split(string(s), "|")
	*str = ss
	return nil
}

func (str StrSlice) Value() (driver.Value, error) {
	strs := strings.Join(str, "|")
	return strs, nil
}
