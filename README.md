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
	DB    string `key:"DBCONN"`
	MAYBE string `key:",optional"`
}

var myConfig = testConf{}

func main() {
	err := dumbconf.Populate(&myConfig)
	if err == nil {
		log.Printf("myConfig = %+v", myConfig)
	}
}
```
```Terminal
> go build test.go
> DBCONN="psql:5432" API="http://api4u.com/do" ./test
2021/07/17 20:56:37 myConfig = {API:http://api4u.com/do DB:psql:5432 MAYBE:}
> ./test -h
usage of: ./test
  -api string
    
  -db string
    
  -maybe string

> DBCONN="psql:5432" ./test -api "http://api4u.com/do" -maybe yes
2021/07/24 00:46:04 myConfig = {API:http://api4u.com/do DB:psql:5432 MAYBE:yes}
```
