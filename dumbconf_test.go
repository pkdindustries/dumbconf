package dumbconf

import (
	"os"
	"testing"
)

type testConf struct {
	DB    string
	API   string `key:"APIBASE"`
	MAYBE string `key:",optional"`
}

func setTestEnv(t *testing.T) {
	os.Setenv("DB", "mysql:3306")
	os.Setenv("APIBASE", "https://api4u.com/")
	os.Setenv("MAYBE", "sdf,kj23-r2")
	t.Cleanup(func() {
		os.Unsetenv("DB")
		os.Unsetenv("APIBASE")
		os.Unsetenv("MAYBE")
	})
}

func TestLoadEnv(t *testing.T) {
	setTestEnv(t)
	conf := testConf{}
	err := populate(&conf, false)
	if err != nil {
		t.Fatalf("error should be nil: %v", err)
	}
	if conf.DB != os.Getenv("DB") {
		t.Fatalf("testConf %v", conf)
	}
	if conf.API != os.Getenv("APIBASE") {
		t.Fatalf("testConf %v", conf)
	}
	if conf.MAYBE != os.Getenv("MAYBE") {
		t.Fatalf("testConf %v", conf)
	}
}

func TestLoadEnvMissingOptionalVar(t *testing.T) {
	setTestEnv(t)
	os.Unsetenv("MAYBE")
	conf := testConf{}
	err := populate(&conf, false)
	if err != nil {
		t.Fatalf("error should be nil: %v", err)
	}
}

func TestLoadEnvMissingRequiredVar(t *testing.T) {
	setTestEnv(t)
	os.Unsetenv("APIBASE")
	conf := testConf{}
	err := populate(&conf, false)
	if err == nil {
		t.Fatalf("unset env but no error")
	} else if err.Error() != "APIBASE" {
		t.Fatalf("wrong error: %v", err)
	}
}
