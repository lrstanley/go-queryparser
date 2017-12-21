# go-queryparser [![godoc](https://godoc.org/github.com/lrstanley/go-queryparser?status.png)](https://godoc.org/github.com/lrstanley/go-queryparser) [![goreport](https://goreportcard.com/badge/github.com/lrstanley/go-queryparser)](https://goreportcard.com/report/github.com/lrstanley/go-queryparser) [![build](https://travis-ci.org/lrstanley/go-queryparser.svg?branch=master)](https://travis-ci.org/lrstanley/go-queryparser) [![codecov](https://codecov.io/gh/lrstanley/go-queryparser/branch/master/graph/badge.svg)](https://codecov.io/gh/lrstanley/go-queryparser)

go-queryparser parses a common "q" http GET variable to strip out filters,
which can be used for advanced searching, like:

```
Hello World tags:example,world foo:"something quoted" author:lrstanley
```

## Use:

This will pull from master (currently `v2`):

```console
$ go get -u -v github.com/lrstanley/go-queryparser
```

### v1

```console
$ go get -u -v gopkg.in/lrstanley/go-queryparser.v1
```

## Example:

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/lrstanley/go-queryparser"
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

## Contributing

Please review the [CONTRIBUTING](https://github.com/lrstanley/go-queryparser/blob/master/CONTRIBUTING.md)
doc for submitting issues/a guide on submitting pull requests and helping out.

## License

    LICENSE: The MIT License (MIT)
    Copyright (c) 2017 Liam Stanley <me@liamstanley.io>

    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in
    all copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE.
