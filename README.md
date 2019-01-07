errgen
==

[![Build Status](https://travis-ci.org/moznion/go-errgen.svg?branch=master)](https://travis-ci.org/moznion/go-errgen)

`errgen` is a code generator. This library generates functions that return error according to the definition described in the `struct`.

And the function can receive variables as the function parameters. The parameters will be bound into placeholders that are compatible with sprintf.

Please refer to the [synopsis](#Synopsis) for a concrete example.

Installation
--

```shell
$ go get -u github.com/moznion/go-errgen/cmd/errgen
```

Usage
--

```
Usage of errgen:
  -out-file string
        [optional] the output destination path of the generated code
  -prefix string
        [optional] prefix of error type (default "ERR-")
  -type string
        [mandatory] struct type name of source of error definition
  -version
        show version and revision
```

Synopsis
--

Define errors and configure with `go:generate`:

```go
package mypkg

//go:generate errgen -type=myErrors
type myErrors struct {
	FooErr error `errmsg:"this is FOO error"`
	BarErr error `errmsg:"this is BAR error [%d, %s]" vars:"hoge int, fuga string"`
}
```

And execute go generate:

```shell
$ go generate ./...
```

Then it generates `my_errors_errmsg_gen.go`. That has the following contents:

```go
// This package was auto generated.
// DO NOT EDIT BY YOUR HAND!

package mypkg

import "errors"
import "fmt"

func FooErr() error {
	return errors.New("[ERR-1] this is FOO error")
}

func BarErr(hoge int, fuga string) error {
	return fmt.Errorf("[ERR-2] this is BAR error [%d, %s]", hoge, fuga)
}

func MyErrorsList() []string {
	return []string{
		`[ERR-1] this is FOO error`,
		`[ERR-2] this is BAR error [%d, %s]`,
	}
}
```

Custom Tag Syntax
--

example:

```go
type myErrors struct {
	FooErr error `errmsg:"this is FOO error [%d, %s]" vars:"hoge int, fuga string"`
	BarErr error `errmsg:"this is BAR error" obsoleted:"true"`
}
```

### `errmsg`

- This is a __mandatory__ parameter
- Generated function returns this value
- This parameter supports `sprintf` style placeholder
  - If you use the placeholder, you should use the `vars` parameter together

### `vars`

- This is an optional parameter
- The generated function uses this value as a function parameter
  - i.e. this parameter must be the valid syntax of golang's function parameter
- And variables that are described by this parameter will be filled into `sprintf` style placeholders of `errmsg`

### `obsoleted`

- This is an optional parameter
- If this parameter is __not empty__, the error message won't be generated as go code
  - But it increments a serial number of prefix
  - This parameter is useful for excluding the error message from the target of code generation
    - See also [Notes](#Notes)

Notes
--

- errgen will automatically give a prefix to each error message
  - The prefix contains a serial number
    - It is useful to determine the error type
    - So you should not...
      - Remove an error message definition from the struct
      - Change order of error messages
    - => __You should only append error definitions__
      - Please consider using `obsoleted` parameter to exclude an error message from the target of code generation

FAQ
--

### Why don't use map like structure instead of struct?

A map doesn't ensure the order of members. That is a bit inconvenient to generate error messages with a prefix that contains the serial number.

### Why don't use slice like structure instead of struct?

It's a good point. But I (i.e. moznion) don't want to consider the unbalanced slice members; unbalanced means odd number of members.

License
--

```
The MIT License (MIT)
Copyright © 2019 moznion, http://moznion.net/ <moznion@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```

