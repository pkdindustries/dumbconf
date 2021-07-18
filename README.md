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
package main
    
import (
	"log"

	"github.com/pkdindustries/dumbconf"
)

type testConf struct {
	API   string
	DB    string `env:"DBCONN"`
	MAYBE string `env:"MAYBE,optional"`
}

var myConfig = testConf{}

func main() {
	err := dumbconf.LoadConfig(&myConfig)
	if err == nil {
		log.Printf("myConfig = %+v", myConfig)
	}
}
```
```Terminal
> export DBCONN="psql:5432"
> export API="http://api4u.com/do
> go build test.go
> ./test
2021/07/17 20:56:37 conf = {API:http://api4u.com/do DB:psql:5432 MAYBE:}
```
