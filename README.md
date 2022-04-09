<!-- template:begin:header -->
<!-- template:end:header -->

<!-- template:begin:toc -->
<!-- template:end:toc -->

## What?

go-queryparser parses a common "q" http GET variable to strip out filters,
which can be used for advanced searching, like:

```
Hello World tags:example,world foo:"something quoted" author:lrstanley
```

## Use:

This will pull from master (currently `v2`):

```console
$ go get -u github.com/lrstanley/go-queryparser@latest
```

## Example:

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/lrstanley/v2/go-queryparser"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		q := queryparser.Parse(r.FormValue("q"))
		if q.Has("author") {
			fmt.Fprintf(w, "filtering by author %q!\n", q.GetOne("author"))
			return
		}

		fmt.Fprint(w, "no filtering requested!\n")
	})

	http.ListenAndServe(":8080", nil)
}
```

```console
$ curl -s localhost:8080
no filtering requested!
$ curl -s 'localhost:8080?q=author:"liam"'
filtering by author "liam"!
```

The main benefit is for user input boxes where you want additional filtering,
like the Github issues search box, or similar.

<!-- template:begin:support -->
<!-- template:end:support -->

<!-- template:begin:contributing -->
<!-- template:end:contributing -->

<!-- template:begin:license -->
<!-- template:end:license -->
