package dumbconf

import (
	"os"
	"testing"
)

type testConf struct {
	DB    string
	API   string `env:"APIBASE"`
	TOKEN string `env:"special_token"`
}

func setTestEnv(t *testing.T) {
	os.Setenv("DB", "mysql:3306")
	os.Setenv("APIBASE", "https://api4u.com/")
	os.Setenv("special_token", "sdf,kj23-r2")
	t.Cleanup(func() {
		os.Unsetenv("DB")
		os.Unsetenv("APIBASE")
		os.Unsetenv("special_token")
	})
}

func TestLoadEnv(t *testing.T) {

	setTestEnv(t)

	conf := testConf{}
	err := LoadConfig(&conf)

	if err != nil {
		t.Failed()
	}
	if conf.DB != os.Getenv("DB") {
		t.Fatalf("TestConfiguration %v", conf)
	}
	if conf.API != os.Getenv("APIBASE") {
		t.Fatalf("TestConfiguration %v", conf)
	}
	if conf.TOKEN != os.Getenv("special_token") {
		t.Fatalf("TestConfiguration %v", conf)
	}
}

func TestLoadEnvFail(t *testing.T) {

	setTestEnv(t)
	os.Unsetenv("APIBASE")

	conf := testConf{}
	err := LoadConfig(&conf)

	if err == nil {
		t.Fatalf("%v", err)
	} else if err.Error() != "APIBASE" {
		t.Fatalf("wrong error: %v", err)
	}
}
