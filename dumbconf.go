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
//		API   string `key:"APIBASE"`
//		MAYBE string `key:"MAYBE,optional"`
//	}
type dumbStruct interface{}

// extracted metadata from the dumb struct
type dumbMap map[string]*configVar
type configVar struct {
	Key       string
	Value     string
	FlagValue string
	Optional  bool
}

// pass in a pointer to struct, have it populated from system environment
// and command line flags, defaulted to field name, or via 'key' tag
//	type Config struct {
//		DB    string
//		API   string `key:"APIBASE"`
//		MAYBE string `key:"MAYBE,optional"`
//	}
// 	var c = Config{}
//	err := dumbconf.Populate(&c)
func Populate(dumb dumbStruct) error {
	return populate(dumb, true)
}

// reflects on the conf struct, parses args and creates a map
// of ConfigVar
func createMap(dumb dumbStruct) dumbMap {
	c := make(dumbMap)
	mirror := reflect.ValueOf(dumb)
	typeof := mirror.Elem().Type()
	for i := 0; i < typeof.NumField(); i++ {
		field := typeof.Field(i)
		cfg := &configVar{}
		env, ok := field.Tag.Lookup("key")
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
func (c dumbMap) unmarshall(dumb dumbStruct) error {
	for fieldname := range c {
		f := reflect.ValueOf(dumb).Elem().FieldByName(fieldname)
		if f.CanSet() {
			cfg := c[fieldname]
			value := cfg.FlagValue
			if value == "" {
				value = cfg.Value
			}
			if value == "" && !cfg.Optional {
				return fmt.Errorf("missing config: %s", c.missingKeys())
			}
			f.SetString(value)
		} else {
			return errors.New("[!]")
		}
	}
	return nil
}

func (c dumbMap) missingKeys() []string {
	ks := []string{}
	for f := range c {
		if !c[f].Optional && c[f].Value == "" && c[f].FlagValue == "" {
			ks = append(ks, c[f].Key)
		}
	}
	return ks
}

// read flag values into config
func (c dumbMap) readFlags() {
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage of: %s\n", os.Args[0])
		flags.PrintDefaults()
	}
	for k := range c {
		cfg := c[k]
		flagname := strings.ToLower(k)
		flags.StringVar(&cfg.FlagValue, flagname, "", "")
	}
	flags.Parse(os.Args[1:])
}

// read environment values into config
func (c dumbMap) readEnv() {
	for k := range c {
		cfg := c[k]
		value, _ := os.LookupEnv(cfg.Key)
		cfg.Value = value
	}
}

// create config, merge flags>env .. push back into
// user's struct
func populate(dumb dumbStruct, flags bool) error {
	cfgmap := createMap(dumb)
	cfgmap.readEnv()
	if flags {
		cfgmap.readFlags()
	}
	err := cfgmap.unmarshall(dumb)
	return err
}
