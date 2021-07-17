<h2> config your golang, with dumb </h2>

```
go get github.com/pkdindustries/dumbconf
```

```
import (
    "log"
    "os"
    "github.com/pkdindustries/dumbconf"
)

type testConf struct {
    DB    string `env:"ENVKEY_DB"`
}

var myConfig = testConf{}

func main() {
    os.Setenv("ENVKEY_DB","psql:5432")
    dumbconf.loadConfig(&myConfig)
    log.Printf("conf = %v", myConfig)
}
```