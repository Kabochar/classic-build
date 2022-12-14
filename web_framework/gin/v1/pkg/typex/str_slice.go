package typex

import (
	"database/sql/driver"
	"strings"
)

// 给gorm提供数组类型的支持(参考案例Scan方式的类型判断有问题)
// https://blog.csdn.net/js010111/article/details/126076320
type StrSlice []string

func (s *StrSlice) Scan(val interface{}) error {
	value := val.(string)
	split := strings.Split(string(value), "|")
	*s = split
	return nil
}

func (s StrSlice) Value() (driver.Value, error) {
	value := strings.Join(s, "|")
	return value, nil
}
