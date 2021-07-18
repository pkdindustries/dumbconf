package dumbconf

import (
	"log"
	"os"
	"reflect"
	"strings"
)

type dumbError struct {
	key string
}

func (d *dumbError) Error() string {
	return d.key
}

// pass in a pointer to struct, have it populated from system environment.
//		type Config struct {
//			DB    string // defaults to env:"DB,required"
//			API   string `env:"APIBASE"`
//			MAYBE string `env:"MAYBE,optional"`
//		}
// 		var c = Config{}
//		err := dumbconf.LoadConfig(&c)
func LoadConfig(conf interface{}) error {
	mirror := reflect.ValueOf(conf)
	typeof := mirror.Elem().Type()
	for i := 0; i < typeof.NumField(); i++ {
		field := typeof.Field(i)
		key := ""
		opt := ""

		env, ok := field.Tag.Lookup("env")
		if !ok {
			key = field.Name
		} else {
			tags := strings.Split(env, ",")
			key = tags[0]
			if key == "" {
				key = field.Name
			}
			if len(tags) > 1 {
				opt = tags[1]
			}
		}

		value, ok := os.LookupEnv(key)
		if !ok {
			log.Printf("dumbconf: [%s] unset", key)
			if opt == "" {
				return &dumbError{key}
			}
		}

		stype := reflect.ValueOf(conf).Elem()
		sfield := stype.Field(i)
		if sfield.CanSet() {
			sfield.SetString(value)
		} else {
			return &dumbError{}
		}
	}
	return nil
}
