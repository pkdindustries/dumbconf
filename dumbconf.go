package dumbconf

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// user supplied struct with tags
//	type DumbConfig struct {
//		DB    string
//		API   string `env:"APIBASE"`
//		MAYBE string `env:"MAYBE,optional"`
//	}
type DumbStruct interface{}

// extracted metadata from the dumb struct
type ConfigMap map[string]*ConfigVar
type ConfigVar struct {
	Key       string
	Value     string
	FlagValue string
	Optional  bool
}

// reflects on the conf struct, parses args and creates a map
func createMap(conf DumbStruct) ConfigMap {
	c := make(ConfigMap)
	mirror := reflect.ValueOf(conf)
	typeof := mirror.Elem().Type()
	for i := 0; i < typeof.NumField(); i++ {
		field := typeof.Field(i)
		cfg := &ConfigVar{}
		env, ok := field.Tag.Lookup("env")
		if !ok {
			cfg.Key = field.Name
		} else {
			args := strings.Split(env, ",")
			key := args[0]
			if key == "" {
				key = field.Name
			}
			cfg.Key = key
			for _, opt := range args[1:] {
				switch opt {
				case "optional":
					cfg.Optional = true
				}
			}
		}
		c[field.Name] = cfg
	}
	return c
}

// merges flag and env values, writes struct
func (c ConfigMap) unmarshall(conf DumbStruct) error {
	for fieldname := range c {
		f := reflect.ValueOf(conf).Elem().FieldByName(fieldname)
		if f.CanSet() {
			cfg := c[fieldname]
			value := cfg.FlagValue
			if value == "" {
				value = cfg.Value
			}
			if value == "" && !cfg.Optional {
				return errors.New(cfg.Key)
			}
			f.SetString(value)
		} else {
			return errors.New("[!]")
		}
	}
	return nil
}

// read flag values into config
func (c ConfigMap) readFlags() {
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage of: %s\n", os.Args[0])
		flags.PrintDefaults()
	}
	for k := range c {
		cfg := c[k]
		flagname := strings.ToLower(cfg.Key)
		flags.StringVar(&cfg.FlagValue, flagname, "", "")
	}
	flags.Parse(os.Args[1:])
}

// read environment values into config
func (c ConfigMap) readEnv() {
	for k := range c {
		cfg := c[k]
		value, _ := os.LookupEnv(cfg.Key)
		cfg.Value = value
	}
}

// pass in a pointer to struct, have it populated from system environment.
//	type Config struct {
//		DB    string
//		API   string `env:"APIBASE"`
//		MAYBE string `env:"MAYBE,optional"`
//	}
// 	var c = Config{}
//	err := dumbconf.Populate(&c)
func populate(dumb DumbStruct, flags bool) error {
	cfgmap := createMap(dumb)
	cfgmap.readEnv()
	if flags {
		cfgmap.readFlags()
	}
	err := cfgmap.unmarshall(dumb)
	return err
}

func Populate(dumb DumbStruct) error {
	return populate(dumb, true)
}
