# tuples

[![codecov](https://codecov.io/gh/antklim/tuples/branch/main/graph/badge.svg?token=8RHXR5OABD)](https://codecov.io/gh/antklim/tuples)

Go package `tuples` reads and writes tuples string - a string containing a list of records/tuples separated by a delimiter.

It's not rare when a configuration parameter represents a record or a list of records. For example:
```
SUPPORTED_FORMATS="height=700,width=350,format=jpeg height=900,width=450,format=png"
```

The `tuples` package allows to unmarshal such strings to a predefined Go structure or an arbitrary map. It also supports structure and map marshaling to a tuples string.

# Installation
`go get github.com/antklim/tuples`

# Format
The tuples string format is the following:
```
<name=value>,... [<name=value>,...]
```

There are three delimiters in format:
* key-value delimiter, the default value is `=`
* filds delimiter, the default value is `,`
* tuples (records) delimiter, the default value is ` ` (a whitespace)

A string can contain 0 to N tuples. Each tuple can consist of 1 to M fields.

# Usage

## Unmarshal

The package uses `tuples` tag followed by the field name to decode to a Go structure. Structure fields without the `tuples` tag omitted during decoding.  
When unmarshaling to a map the tuples string field names become keys in the map. 

The package does not read the full tuples string for decoding. It scans the string tuple by tuple. It is not possible to know ahead how many tuples the string contains. Therefore, the package only accepts the following unmarshaling destinations:
* a slice or array of a struct
* an interface.

In case when an interface (`any`) provided as the decoding destination, a slice of the arbitrary maps produced: `[]map[string]any`. Note, map keys will be alphabetically sorted.

```go
package main

import(
	"fmt"
	"time"

	"github.com/antklim/tuples"
)

type format struct {
	Height int    `tuples:"h"`
	Width  int    `tuples:"w"`
	Format string `tuples:"f"`
}

func main() {
	formatsConf := "h=700,w=350,f=jpeg h=900,w=450,f=png"

	var formats []format
	if err := tuples.Unmarshal([]byte(formatsConf), &formats); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", formats)

	var anys any
	if err := tuples.Unmarshal([]byte("b=2,a=1 d=4,c=3"), &anys); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", anys)

	// Output:
	// [{Height:700 Width:350 Format:jpeg} {Height:900 Width:450 Format:png}]
	// [map[a:1 b:2] map[c:3 d:4]]
}
```

Additionally, the package provides a `Reader`. It reads a tuples string and produces a collection of tuple values. You can read all tuples at once:

```go
package main

import(
	"fmt"
	"strings"
	"time"

	"github.com/antklim/tuples"
)

func main() {
	formatsConf := "h=700,w=350,f=jpeg h=900,w=450,f=png"

	r, err := tuples.NewReader(strings.NewReader(formatsConf))
	if err != nil {
		fmt.Println(err)
	}

	v, err := r.ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", v)

	// Output:
	// [[700 350 jpeg] [900 450 png]]
}
```

Or read tuple by tuple.

```go
package main

import(
	"fmt"
	"strings"
	"time"

	"github.com/antklim/tuples"
)

func main() {
	formatsConf := "h=700,w=350,f=jpeg h=900,w=450,f=png"

	r, err := tuples.NewReader(strings.NewReader(formatsConf))
	if err != nil {
		fmt.Println(err)
	}

	for {
		tuple, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%v\n", tuple)
	}

	// Output:
	// [700 350 jpeg]
 	// [900 450 png]
}
```

## Marshal
