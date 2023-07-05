<!-- template:define:options
{
  "nodescription": true
}
-->
![logo](https://liam.sh/-/gh/svg/lrstanley/go-queryparser?icon=mdi%3Afilter-cog-outline&icon.height=110&layout=left&icon.color=rgba%280%2C+184%2C+126%2C+1%29)
<!-- template:begin:header -->
<!-- do not edit anything in this "template" block, its auto-generated -->

<p align="center">
  <a href="https://github.com/lrstanley/go-queryparser/tags">
    <img title="Latest Semver Tag" src="https://img.shields.io/github/v/tag/lrstanley/go-queryparser?style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/go-queryparser/commits/master">
    <img title="Last commit" src="https://img.shields.io/github/last-commit/lrstanley/go-queryparser?style=flat-square">
  </a>


  <a href="https://github.com/lrstanley/go-queryparser/actions?query=workflow%3Atest+event%3Apush">
    <img title="GitHub Workflow Status (test @ master)" src="https://img.shields.io/github/actions/workflow/status/lrstanley/go-queryparser/test.yml?branch=master&label=test&style=flat-square">
  </a>

  <a href="https://codecov.io/gh/lrstanley/go-queryparser">
    <img title="Code Coverage" src="https://img.shields.io/codecov/c/github/lrstanley/go-queryparser/master?style=flat-square">
  </a>

  <a href="https://pkg.go.dev/github.com/lrstanley/go-queryparser/v3">
    <img title="Go Documentation" src="https://pkg.go.dev/badge/github.com/lrstanley/go-queryparser/v3?style=flat-square">
  </a>
  <a href="https://goreportcard.com/report/github.com/lrstanley/go-queryparser/v3">
    <img title="Go Report Card" src="https://goreportcard.com/badge/github.com/lrstanley/go-queryparser/v3?style=flat-square">
  </a>
</p>
<p align="center">
  <a href="https://github.com/lrstanley/go-queryparser/issues?q=is:open+is:issue+label:bug">
    <img title="Bug reports" src="https://img.shields.io/github/issues/lrstanley/go-queryparser/bug?label=issues&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/go-queryparser/issues?q=is:open+is:issue+label:enhancement">
    <img title="Feature requests" src="https://img.shields.io/github/issues/lrstanley/go-queryparser/enhancement?label=feature%20requests&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/go-queryparser/pulls">
    <img title="Open Pull Requests" src="https://img.shields.io/github/issues-pr/lrstanley/go-queryparser?label=prs&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/go-queryparser/discussions/new?category=q-a">
    <img title="Ask a Question" src="https://img.shields.io/badge/support-ask_a_question!-blue?style=flat-square">
  </a>
  <a href="https://liam.sh/chat"><img src="https://img.shields.io/badge/discord-bytecord-blue.svg?style=flat-square" title="Discord Chat"></a>
</p>
<!-- template:end:header -->

<!-- template:begin:toc -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :link: Table of Contents

  - [What?](#what)
  - [Use:](#use)
  - [Example:](#example)
  - [Support &amp; Assistance](#raising_hand_man-support--assistance)
  - [Contributing](#handshake-contributing)
  - [⚖️ License](#balance_scale-license)
<!-- template:end:toc -->

## What?

go-queryparser parses a common "q" http GET variable to strip out filters,
which can be used for advanced searching, like:

```
Hello World tags:example,world foo:"something quoted" author:lrstanley
```

## Use:

<!-- template:begin:goget -->
<!-- do not edit anything in this "template" block, its auto-generated -->
```console
go get -u github.com/lrstanley/go-queryparser/v3@latest
```
<!-- template:end:goget -->

## Example:

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/lrstanley/go-queryparser/v3"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		q := queryparser.Parse(r.FormValue("q"))
		if q.Has("author") {
			fmt.Fprintf(w, "filtering by author %q!", q.GetOne("author"))
			return
		}

		fmt.Fprint(w, "no filtering requested!")
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
<!-- do not edit anything in this "template" block, its auto-generated -->
## :raising_hand_man: Support & Assistance

* :heart: Please review the [Code of Conduct](.github/CODE_OF_CONDUCT.md) for
     guidelines on ensuring everyone has the best experience interacting with
     the community.
* :raising_hand_man: Take a look at the [support](.github/SUPPORT.md) document on
     guidelines for tips on how to ask the right questions.
* :lady_beetle: For all features/bugs/issues/questions/etc, [head over here](https://github.com/lrstanley/go-queryparser/issues/new/choose).
<!-- template:end:support -->

<!-- template:begin:contributing -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :handshake: Contributing

* :heart: Please review the [Code of Conduct](.github/CODE_OF_CONDUCT.md) for guidelines
     on ensuring everyone has the best experience interacting with the
    community.
* :clipboard: Please review the [contributing](.github/CONTRIBUTING.md) doc for submitting
     issues/a guide on submitting pull requests and helping out.
* :old_key: For anything security related, please review this repositories [security policy](https://github.com/lrstanley/go-queryparser/security/policy).
<!-- template:end:contributing -->

<!-- template:begin:license -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :balance_scale: License

```
MIT License

Copyright (c) 2017 Liam Stanley <me@liamstanley.io>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

_Also located [here](LICENSE)_
<!-- template:end:license -->
