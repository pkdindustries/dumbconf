<h2> If you try harder than this to load config from the environment, you should stop. You are only hurting yourself. </h2>
<br>
<h3> features </h3>
<li> strings! 
<li> env key to struct field mapping
<li> optional fields
<p>
<h3> non-features </h3>
<li> types other than strings </li>
<li> defaults </li>
<p>
<br>

```Terminal
> go get github.com/pkdindustries/dumbconf
```

```Go
import (
    "log"
    "os"
    "github.com/pkdindustries/dumbconf"
)

type testConf struct {
    API     string  
    DB      string `env:"ENVKEY_DB"`
    MAYBE   string `env:"MAYBE,optional`
}

var myConfig = testConf{}

func main() {
    os.Setenv("ENVKEY_DB","psql:5432")
    os.Setenv("API","https://api4u.com/do")
    os.Unsetenv("MAYBE")
    dumbconf.LoadConfig(&myConfig)
    log.Printf("conf = %v", myConfig)
}
```
