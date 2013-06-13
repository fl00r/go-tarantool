# Tarantool

[Tarantool](http://tarantool.org/) client on Go.

## Usage

```go
package main

import "github.com/fl00r/go-tarantool"

func main() {
	connection = tarantool.Connect("loaclhost:33013")
	space = connection.Space(1)
	// space.Insert()
	// space.Select()
	// space.Delete()
	// space.Update()
}
```

TO BE WITTEN