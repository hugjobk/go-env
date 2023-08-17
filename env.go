package env

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"syscall"

	"github.com/hugjobk/go-text"
)

func GetEnv(k string, v interface{}) (bool, error) {
	e, ok := syscall.Getenv(k)
	if !ok {
		return false, nil
	}
	if err := text.Unmarshal(e, v); err != nil {
		return true, err
	}
	return true, nil
}

func ParseEnv(v interface{}) error {
	tp := reflect.TypeOf(v)
	if tp.Kind() != reflect.Ptr || tp.Elem().Kind() != reflect.Struct {
		return errors.New("v must be pointer to a struct")
	}
	t := tp.Elem()
	rv := reflect.Indirect(reflect.ValueOf(v))
	for i := 0; i < rv.NumField(); i++ {
		f := rv.Field(i)
		if !f.CanSet() {
			continue
		}
		if f.Kind() == reflect.Struct {
			if err := ParseEnv(f.Addr().Interface()); err != nil {
				return err
			}
		}
		tag, ok := t.Field(i).Tag.Lookup("env")
		if !ok {
			continue
		}
		if !strings.Contains(tag, ",") {
			if _, err := GetEnv(tag, f.Addr().Interface()); err != nil {
				return fmt.Errorf("env.ParseEnv: %s: %s", tag, err)
			}
		} else {
			s := strings.SplitN(tag, ",", 2)
			ok, err := GetEnv(s[0], f.Addr().Interface())
			if err != nil {
				return fmt.Errorf("env.ParseEnv: %s: %s", s[0], err)
			}
			if !ok {
				if err := text.Unmarshal(s[1], f.Addr().Interface()); err != nil {
					return fmt.Errorf("env.ParseEnv: %s: %s", s[0], err)
				}
			}
		}

	}
	return nil
}
