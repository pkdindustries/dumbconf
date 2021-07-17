package dumbconf

import (
	"log"
	"os"
	"reflect"
)

type dumbError struct {
	key string
}

func (d *dumbError) Error() string {
	return d.key
}

func LoadConfig(conf interface{}) error {
	mirror := reflect.ValueOf(conf)
	typeof := mirror.Elem().Type()
	for i := 0; i < typeof.NumField(); i++ {
		field := typeof.Field(i)
		key, ok := field.Tag.Lookup("env")
		if !ok {
			key = field.Name
		}
		value := os.Getenv(key)
		if value == "" {
			log.Printf("key '%s' not found in environment", key)
			return &dumbError{key}
		} else {
			log.Printf("key '%s' found in environment", key)
		}
		stype := reflect.ValueOf(conf).Elem()
		sfield := stype.Field(i)
		if sfield.CanSet() {
			sfield.SetString(value)
		} else {
			return &dumbError{key}
		}
	}
	return nil
}
